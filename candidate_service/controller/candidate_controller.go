package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"example.com/se/config"
	"example.com/se/entity"
)

// struct สำหรับ response รวมข้อมูล election (ถ้ารู้โครงสร้าง election ให้แทน interface{} ด้วย struct จริง)
type CandidateWithElection struct {
	ID         uint                   `json:"id"`
	Name       string                 `json:"name"`
	ElectionID uint                   `json:"election_id"`
	Election   map[string]interface{} `json:"election"` // ใช้ map เพื่อให้แปลง json dynamic ได้ง่าย
}

// ฟังก์ชันดึงข้อมูล election จาก election_service
func fetchElection(electionID uint) (map[string]interface{}, error) {
	url := fmt.Sprintf("http://election_service:8002/api/election/%d", electionID)
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Failed to call election service: %v\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errMsg := fmt.Errorf("failed to get election data, status code: %d", resp.StatusCode)
		log.Println(errMsg)
		return nil, errMsg
	}

	var electionData map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&electionData); err != nil {
		log.Printf("Failed to decode election data: %v\n", err)
		return nil, err
	}
	return electionData, nil
}

// GetAllCandidates ดึงข้อมูล candidates พร้อม election (เรียกทีละตัว)
func GetAllCandidates(c *gin.Context) {
	db := config.DB()
	var candidates []entity.Candidates

	if err := db.Find(&candidates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var result []CandidateWithElection

	for _, candidate := range candidates {
		electionData, err := fetchElection(candidate.ElectionID)
		if err != nil {
			log.Printf("Error fetching election for election_id=%d: %v", candidate.ElectionID, err)
			electionData = nil
		} else {
			log.Printf("Election data for election_id=%d: %v", candidate.ElectionID, electionData)
		}

		result = append(result, CandidateWithElection{
			ID:         candidate.ID,
			Name:       candidate.Name,
			ElectionID: candidate.ElectionID,
			Election:   electionData,
		})
	}


	c.JSON(http.StatusOK, result)
}


// func GetAllCandidates(c *gin.Context) {
// 	db := config.DB()
// 	var candidates []entity.Candidates

// 	if err := db.Find(&candidates).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	fmt.Printf("Candidates found: %v\n", candidates)

// 	c.JSON(http.StatusOK, candidates)  // ส่งตรงข้อมูล candidates ไปก่อน
// }

// GetCandidate ดึง candidate เดี่ยวพร้อม election
func GetCandidate(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid candidate ID"})
		return
	}

	db := config.DB()
	var candidate entity.Candidates
	if err := db.First(&candidate, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Candidate not found"})
		return
	}

	electionData, err := fetchElection(candidate.ElectionID)
	if err != nil {
		// แจ้งเตือนแต่ยังส่ง candidate กลับไป
		c.JSON(http.StatusOK, gin.H{
			"candidate": candidate,
			"warning":   "Failed to fetch election data",
		})
		return
	}

	response := CandidateWithElection{
		ID:         candidate.ID,
		Name:       candidate.Name,
		ElectionID: candidate.ElectionID,
		Election:   electionData,
	}

	c.JSON(http.StatusOK, response)
}

func CreateCandidate(c *gin.Context) {
	var candidate entity.Candidates
	if err := c.ShouldBindJSON(&candidate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := config.DB()
	if err := db.Create(&candidate).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, candidate)
}

func UpdateCandidate(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid candidate ID"})
		return
	}

	var candidate entity.Candidates
	if err := c.ShouldBindJSON(&candidate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := config.DB()
	var existing entity.Candidates
	if err := db.First(&existing, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Candidate not found"})
		return
	}

	candidate.ID = existing.ID
	if err := db.Save(&candidate).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, candidate)
}

func DeleteCandidate(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid candidate ID"})
		return
	}

	db := config.DB()
	if err := db.Delete(&entity.Candidates{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Candidate deleted"})
}
