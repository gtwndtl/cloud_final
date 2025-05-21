package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"example.com/se/config"
	"example.com/se/entity"
	"example.com/se/metrics"
)

// Candidate struct สำหรับดึงข้อมูลจาก candidate_service
type Candidate struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	ElectionID uint   `json:"election_id"`
}

// ดึง Elections ทั้งหมด พร้อม candidate_ids
func GetAllElections(c *gin.Context) {
	db := config.DB()
	var elections []entity.Elections
	if err := db.Find(&elections).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// สำหรับแต่ละ election ให้ดึง candidate_ids จาก candidate_service
	for i := range elections {
		candidates, err := fetchCandidatesByElectionID(elections[i].ID)
		if err != nil {
			// log error, แต่ยังส่งข้อมูล election กลับได้ (candidate_ids จะว่าง)
			elections[i].CandidateIDs = []uint{}
			continue
		}

		var candidateIDs []uint
		for _, candidate := range candidates {
			candidateIDs = append(candidateIDs, candidate.ID)
		}
		elections[i].CandidateIDs = candidateIDs
	}

	c.JSON(http.StatusOK, elections)
}

// ดึง election ทีละตัว พร้อม candidate_ids
func GetElection(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid election ID"})
		return
	}

	db := config.DB()
	var election entity.Elections
	if err := db.First(&election, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Election not found"})
		return
	}

	candidates, err := fetchCandidatesByElectionID(election.ID)
	if err != nil {
		election.CandidateIDs = []uint{}
	} else {
		var candidateIDs []uint
		for _, candidate := range candidates {
			candidateIDs = append(candidateIDs, candidate.ID)
		}
		election.CandidateIDs = candidateIDs
	}

	c.JSON(http.StatusOK, election)
}

func CreateElection(c *gin.Context) {
	var election entity.Elections
	if err := c.ShouldBindJSON(&election); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if election.Status == "" {
		election.Status = "pending"
	}
	if election.StartTime.IsZero() {
		election.StartTime = time.Now()
	}

	db := config.DB()
	if err := db.Create(&election).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// เพิ่ม counter create
	metrics.ElectionsCreateTotal.Inc()

	// อัปเดต gauge จำนวน election ทั้งหมด
	var count int64
	db.Model(&entity.Elections{}).Count(&count)
	metrics.ElectionsTotal.Set(float64(count))

	c.JSON(http.StatusCreated, election)
}

func UpdateElection(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid election ID"})
		return
	}

	var election entity.Elections
	if err := c.ShouldBindJSON(&election); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := config.DB()
	var existing entity.Elections
	if err := db.First(&existing, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Election not found"})
		return
	}

	election.ID = existing.ID
	if err := db.Save(&election).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// เพิ่ม counter update
	metrics.ElectionsUpdateTotal.Inc()

	c.JSON(http.StatusOK, election)
}

func DeleteElection(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid election ID"})
		return
	}

	db := config.DB()
	if err := db.Delete(&entity.Elections{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// เพิ่ม counter delete
	metrics.ElectionsDeleteTotal.Inc()

	// อัปเดต gauge จำนวน election หลังลบ
	var count int64
	db.Model(&entity.Elections{}).Count(&count)
	metrics.ElectionsTotal.Set(float64(count))

	c.JSON(http.StatusOK, gin.H{"message": "Election deleted"})
}


// ฟังก์ชันช่วยดึง candidates จาก candidate_service ผ่าน HTTP API
func fetchCandidatesByElectionID(electionID uint) ([]Candidate, error) {
	url := fmt.Sprintf("http://candidate_service:8003/api/candidates?election_id=%d", electionID)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("candidate service returned status %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var candidates []Candidate
	if err := json.Unmarshal(body, &candidates); err != nil {
		return nil, err
	}

	return candidates, nil
}
