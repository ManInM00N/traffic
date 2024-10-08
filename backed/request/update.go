package request

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"main/binary"
	"net/http"
)

type QueryStationInfo struct {
	ID   uint    ` json:"id" binding:"required"`
	Lon  *string ` json:"lon" binding:"required"`
	Lat  *string ` json:"lat" binding:"required"`
	Name *string ` json:"name" binding:"required"`
	Use  *int    ` json:"use" binding:"required"`
}

func UpdateStation(c *gin.Context) {
	val, has := c.Get("isadmin")
	if !has {
		c.JSON(404, gin.H{
			"error": "system error",
			"code":  0,
		})
		c.Abort()
		return
	}
	can := val.(bool)

	var tmp SubwayStation
	var data QueryStationInfo
	err := c.BindJSON(&data)
	if err != nil {
		c.JSON(417, gin.H{
			"code":  0,
			"error": "bad request",
		})
		c.Abort()
		return
	}
	tmp.ID = data.ID
	tmp.Name = *data.Name
	tmp.Lon = *data.Lon
	tmp.Lat = *data.Lat
	use := *data.Use
	if !can && use == 0 {
		c.JSON(http.StatusForbidden, gin.H{
			"code":  0,
			"error": "Insufficient authority",
		})
		c.Abort()
		return
	}
	val, has = c.Get("staff_id")
	who := val.(string)
	var station SubwayStation
	if use == -1 {
		station.ID = tmp.ID
		if !can {
			var sub Submission
			Db.Last(&sub)
			sub.Done = false
			sub.ID++
			sub.AdminId = ""
			sub.UserId = who
			sub.Message = ""
			sub.Op = -1
			data, err := json.Marshal(tmp)
			if err != nil {
				log.Fatal(err)
			}
			sub.Data = string(data)
			Db.Create(&sub)
			c.JSON(200, gin.H{
				"code": 1,
			})
			return
		}
		if err := Db.Where("id = ?", station.ID).First(&station).Error; err != nil {
			c.JSON(417, gin.H{
				"code":  0,
				"error": "station not found",
			})
			c.Abort()
			return
		}
		Db.Where("subway_station_id = ? ", station.ID).Delete(SubwayStationSubwayline{})
		Db.Model(&SubwayStationSubwayline{}).Where("up = ?", station.ID).Update("up", 0)
		Db.Model(&SubwayStationSubwayline{}).Where("down = ?", station.ID).Update("down", 0)

		Db.Where("id = ?", station.ID).Delete(&station)
		if binary.Setting.Debug {
			binary.DebugLog.Println(station)
		}

	} else if use == 1 {
		if !can {
			var sub Submission
			Db.Last(&sub)
			sub.Done = false
			sub.ID++
			sub.AdminId = ""
			sub.UserId = who
			sub.Op = 1
			sub.Message = ""
			data, err := json.Marshal(tmp)
			if err != nil {
				log.Fatal(err)
			}
			sub.Data = string(data)
			Db.Create(&sub)
			c.JSON(200, gin.H{
				"code": 1,
			})
			return
		}
		Db.Order("id desc").First(&station)
		station.ID++
		station.Name = tmp.Name
		station.Lat = tmp.Lat
		station.Lon = tmp.Lon
		Db.Create(&station)
	} else if use != 0 {
		c.JSON(417, gin.H{
			"code":  0,
			"error": "bad request",
		})
		c.Abort()
		return
	} else {
		if err := Db.Where("id = ? and id != 0", tmp.ID).First(&station).Error; err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"code":  0,
				"error": err,
			})
			c.Abort()
			return
		}
		Db.Where("id = ?", tmp.ID).Save(&tmp)
	}

	var stations []SubwayStation
	Db.Where("id != 0").Find(&stations)
	c.JSON(200, gin.H{
		"code":     1,
		"num":      len(stations),
		"stations": stations,
	})
	return
}

func DealSubmission(c *gin.Context) {
	val, has := c.Get("isadmin")
	if !has {
		c.JSON(404, gin.H{
			"error": "system error",
			"code":  0,
		})
		c.Abort()
		return
	}
	can := val.(bool)
	if !can {
		c.JSON(http.StatusForbidden, gin.H{
			"code":  0,
			"error": "Insufficient authority",
		})
		c.Abort()
		return
	}
	val, has = c.Get("staff_id")
	who := val.(string)
	var sub Submission
	var data map[string]interface{}
	err := c.BindJSON(&data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  0,
			"error": "System error",
		})
		c.Abort()
		return
	}
	id := uint(int(data["id"].(float64)))
	sub.ID = id
	err = Db.Where("id = ?", sub.ID).First(&sub).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  0,
			"error": "Submission not found",
		})
		c.Abort()
		return
	}
	sub.AdminId = who
	if sub.Done {
		c.JSON(http.StatusOK, gin.H{
			"code":  1,
			"error": "Submission already completed",
		})
		c.Abort()
		return
	}

	sub.Done = true
	datas := []byte(sub.Data)
	var next SubwayStation
	err = json.Unmarshal(datas, &next)
	if err != nil {
		sub.Message = "Invalid submission"
		Db.Save(&sub)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":  0,
			"error": "system error",
		})
		return
	}
	if sub.Op == 1 {
		var tmp SubwayStation
		Db.Order("id desc").First(&tmp)
		next.ID = tmp.ID + 1
		Db.Save(&next)
		sub.Message = "Accepted"
	} else {
		var tmp SubwayStation
		err = Db.Where("id = ?", next.ID).First(&tmp).Error
		if err != nil {
			sub.Message = "Station already deleted or not exist"
		} else {
			Db.Delete(&tmp)
			sub.Message = "Accepted"
		}
	}
	Db.Save(&sub)
	c.JSON(200, gin.H{
		"code": 1,
	})
}

type QueryLineInfo struct {
	LineId uint    ` json:"line_id" binding:"required"`
	Name   *string ` json:"name" binding:"required"`
	Use    *int    ` json:"use" binding:"required"`
}

func UpdateLine(c *gin.Context) {
	val, has := c.Get("isadmin")
	if !has {
		c.JSON(404, gin.H{
			"error": "system error",
			"code":  0,
		})
		c.Abort()
		return
	}
	can := val.(bool)
	if !can {
		c.JSON(http.StatusForbidden, gin.H{
			"code":  0,
			"error": "Insufficient authority",
		})
		c.Abort()
		return
	}
	var tmp SubwayLine
	var data QueryLineInfo
	err := c.BindJSON(&data)
	if err != nil {
		c.JSON(417, gin.H{
			"code":  0,
			"error": "failed to bind json",
		})
		c.Abort()
		return
	}
	tmp.LineId = data.LineId
	tmp.Name = *data.Name
	if binary.Setting.Debug {
		binary.DebugLog.Println(tmp.Name)
	}
	use := *data.Use
	var line SubwayLine
	if use == -1 {
		if err := Db.Where("line_id = ?", tmp.LineId).First(&line).Error; err != nil {
			c.JSON(417, gin.H{
				"code":  0,
				"error": "station not found",
			})
			c.Abort()
			return
		}
		Db.Where("subway_line_id = ?", tmp.LineId).Delete(&SubwayStationSubwayline{})
		Db.Where("line_id = ?", tmp.LineId).Delete(&SubwayLine{})

	} else if use == 1 {
		Db.Order("line_id desc").First(&line)
		line.LineId++
		line.Name = tmp.Name
		Db.Create(&line)

	} else if use != 0 {
		c.JSON(417, gin.H{
			"code":  0,
			"error": "bad request",
		})
		c.Abort()
		return
	} else {
		if err = Db.Where("line_id = ? and line_id != 0", tmp.LineId).First(&line).Error; err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"code":  0,
				"error": err,
			})
			c.Abort()
			return
		}
		Db.Where("line_id = ?", tmp.LineId).Save(&tmp)
	}
	var lines []SubwayLine
	Db.Where("line_id != 0").Find(&lines)
	c.JSON(200, gin.H{
		"code":  1,
		"num":   len(lines),
		"lines": lines,
	})
	return
}

type QueryRelationship struct {
	SubwayLineId    uint `json:"line_id" binding:"required"`
	SubwayStationId uint `json:"station_id" binding:"required"`

	Up   *uint ` json:"up" binding:"required"`
	Down *uint ` json:"down" binding:"required"`
	Use  *int  ` json:"use" binding:"required"`
}

func UpdateRelationship(c *gin.Context) {

	val, has := c.Get("isadmin")
	if !has {
		c.JSON(404, gin.H{
			"error": "system error",
			"code":  0,
		})
		c.Abort()
		return
	}
	can := val.(bool)
	if !can {
		c.JSON(http.StatusForbidden, gin.H{
			"code":  0,
			"error": "Insufficient authority",
		})
		c.Abort()
		return
	}
	var data QueryRelationship
	err := c.BindJSON(&data)

	if err != nil {
		c.JSON(417, gin.H{
			"code":  0,
			"error": "bad request",
		})
		binary.DebugLog.Println(data, err)

		c.Abort()
		return
	}
	if binary.Setting.Debug {
		binary.DebugLog.Println(data)
	}
	var tmp SubwayStationSubwayline
	if data.SubwayLineId == 0 || data.SubwayStationId == 0 {
		c.JSON(417, gin.H{
			"code":  0,
			"error": "bad request",
		})
		c.Abort()
		return
	}
	tmp.SubwayLineId = data.SubwayLineId
	tmp.SubwayStationId = data.SubwayStationId
	tmp.Up = *data.Up
	tmp.Down = *data.Down
	use := *data.Use
	if use < -1 || use > 1 {
		c.JSON(417, gin.H{
			"code":  0,
			"error": "bad request",
		})
		c.Abort()
		return
	}

	var num int64
	Db.Model(&SubwayLine{}).Where("line_id = ?", tmp.SubwayLineId).Count(&num)
	if num == 0 {
		c.JSON(417, gin.H{
			"code":  0,
			"error": "subway line not found",
		})
		c.Abort()
		return
	}
	Db.Model(&SubwayStation{}).Where("id = ?", tmp.SubwayStationId).Count(&num)
	if num == 0 {
		c.JSON(417, gin.H{
			"code":  0,
			"error": "subway station not found",
		})
		c.Abort()
		return
	}
	if use == 1 {
		Db.Model(&SubwayStation{}).Where("id = ?", tmp.Up).Count(&num)
		if num == 0 {
			c.JSON(417, gin.H{
				"code":  0,
				"error": "subway station not found",
			})
			c.Abort()
			return
		}
		Db.Model(&SubwayStation{}).Where("id = ?", tmp.Down).Count(&num)
		if num == 0 {
			c.JSON(417, gin.H{
				"code":  0,
				"error": "subway station not found",
			})
			c.Abort()
			return
		}
		Db.Model(&SubwayStationSubwayline{}).Where("subway_station_id = ? and subway_line_id = ?", tmp.SubwayStationId, tmp.SubwayLineId).Count(&num)
		if num == 1 {
			c.JSON(417, gin.H{
				"code":  0,
				"error": "relation already exists",
			})
			c.Abort()
			return
		}
		Db.Create(&tmp)
	} else if use == -1 {
		Db.Model(&SubwayStationSubwayline{}).Where("subway_station_id = ? and subway_line_id = ?", tmp.SubwayStationId, tmp.SubwayLineId).Count(&num)
		if num == 0 {
			c.JSON(417, gin.H{
				"code":  0,
				"error": "relation not found",
			})
			c.Abort()
			return
		}
		Db.Where("subway_station_id = ? and subway_line_id = ?", tmp.SubwayStationId, tmp.SubwayLineId).Delete(&tmp)
	} else {
		Db.Model(&SubwayStation{}).Where("id = ?", tmp.Up).Count(&num)
		if num == 0 {
			c.JSON(417, gin.H{
				"code":  0,
				"error": "subway station not found",
			})
			c.Abort()
			return
		}
		Db.Model(&SubwayStation{}).Where("id = ?", tmp.Down).Count(&num)
		if num == 0 {
			c.JSON(417, gin.H{
				"code":  0,
				"error": "subway station not found",
			})
			c.Abort()
			return
		}
		Db.Model(&SubwayStationSubwayline{}).Where("subway_station_id = ? and subway_line_id = ?", tmp.SubwayStationId, tmp.SubwayLineId).Count(&num)
		if num == 0 {
			c.JSON(417, gin.H{
				"code":  0,
				"error": "relation not found",
			})
			c.Abort()
			return
		}
		Db.Where("subway_station_id = ? and subway_line_id = ?", tmp.SubwayStationId, tmp.SubwayLineId).Save(&tmp)
	}
	var result []SubwayStationSubwayline
	Db.Find(&result)
	c.JSON(200, gin.H{
		"code":   1,
		"num":    len(result),
		"result": result,
	})
}

type QueryTrain struct {
	ID     *string `json:"id" binding:"required"`
	Cap    uint    `json:"cap" binding:"required"`
	LineId *uint   `json:"line_id" binding:"required"`
	Use    *int    `json:"use" binding:"required"`
}

func UpdateTrains(c *gin.Context) {
	val, has := c.Get("isadmin")
	if !has {
		c.JSON(404, gin.H{
			"error": "system error",
			"code":  0,
		})
		c.Abort()
		return
	}
	can := val.(bool)
	if !can {
		c.JSON(http.StatusForbidden, gin.H{
			"code":  0,
			"error": "Insufficient authority",
		})
		c.Abort()
		return
	}
	var data QueryTrain
	err := c.BindJSON(&data)
	if err != nil {
		c.JSON(417, gin.H{
			"code":  0,
			"error": "failed to bind json",
		})
		c.Abort()
		return
	}
	if *data.Use < -1 || *data.Use > 1 {
		c.JSON(400, gin.H{
			"code":  0,
			"error": "bad request",
		})
		c.Abort()
		return
	}
	var num int64
	Db.Model(&SubwayLine{}).Where("line_id = ?", *data.LineId).Count(&num)
	if num == 0 {
		c.JSON(417, gin.H{
			"code":  0,
			"error": "line not found",
		})
		c.Abort()
		return
	}
	Db.Model(&Train{}).Where("id = ?", *data.ID).Count(&num)
	if *data.Use == 1 {
		if num > 0 {
			c.JSON(417, gin.H{
				"code":  0,
				"error": "Train exists",
			})
			c.Abort()
			return
		}
		var tmp Train
		tmp.ID = *data.ID
		tmp.LineId = *data.LineId
		tmp.Capacity = data.Cap
		Db.Create(&tmp)
	} else if *data.Use == 0 {
		if num == 0 {
			c.JSON(417, gin.H{
				"code":  0,
				"error": "Train not found",
			})
			c.Abort()
			return
		}
		Db.Model(&Train{}).Where("id = ?", *data.ID).Updates(Train{ID: *data.ID, Capacity: data.Cap, LineId: *data.LineId})
	} else {
		if num == 0 {
			c.JSON(417, gin.H{
				"code":  0,
				"error": "Train not found",
			})
			c.Abort()
			return
		}
		Db.Model(&Account{}).Where("train_id = ?", data.ID).Update("train_id", "无")
		Db.Model(&Train{}).Where("id = ?", *data.ID).Delete(&Train{})
	}
	var result []Train
	Db.Model(&Train{}).Where("id != '无'").Find(&result)
	c.JSON(200, gin.H{
		"code":   1,
		"num":    len(result),
		"trains": result,
	})
}
