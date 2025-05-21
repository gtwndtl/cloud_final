package controller

import (
    "net/http"
    "strconv"
    "log"
    "encoding/json"
    "fmt"
    "io/ioutil"

    "github.com/gin-gonic/gin"
    "example.com/se/config"
    "example.com/se/entity"
    "example.com/se/metrics"
)

// struct สำหรับตอบกลับ Candidate พร้อม election ข้อมูล
type CandidateWithElection struct {
    ID         uint                   `json:"id"`
    Name       string                 `json:"name"`
    ElectionID uint                   `json:"election_id"`
    Election   map[string]interface{} `json:"election"`
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

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Printf("Failed to read response body: %v\n", err)
        return nil, err
    }

    var electionData map[string]interface{}
    if err := json.Unmarshal(body, &electionData); err != nil {
        log.Printf("Failed to decode election data: %v\n", err)
        return nil, err
    }
    return electionData, nil
}

func GetAllCandidates(c *gin.Context) {
    db := config.DB()
    var candidates []entity.Candidates

    if err := db.Find(&candidates).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // อัปเดต metric candidates total
    metrics.CandidatesTotal.Set(float64(len(candidates)))

    var result []CandidateWithElection

    for _, candidate := range candidates {
        electionData, err := fetchElection(candidate.ElectionID)
        if err != nil {
            log.Printf("Error fetching election for election_id=%d: %v", candidate.ElectionID, err)
            electionData = nil
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

    // เพิ่มจำนวน create metric
    metrics.CandidatesCreateTotal.Inc()

    // อัปเดตจำนวน candidates total
    var count int64
    db.Model(&entity.Candidates{}).Count(&count)
    metrics.CandidatesTotal.Set(float64(count))

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

    // เพิ่มจำนวน update metric
    metrics.CandidatesUpdateTotal.Inc()

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

    // เพิ่มจำนวน delete metric
    metrics.CandidatesDeleteTotal.Inc()

    // อัปเดตจำนวน candidates total หลังลบ
    var count int64
    db.Model(&entity.Candidates{}).Count(&count)
    metrics.CandidatesTotal.Set(float64(count))

    c.JSON(http.StatusOK, gin.H{"message": "Candidate deleted"})
}

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

    // เรียก fetchElection เพื่อดึง election data
    electionData, err := fetchElection(candidate.ElectionID)
    if err != nil {
        // ถ้าดึง election ไม่ได้ ก็แจ้ง warning แต่ส่ง candidate กลับไป
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
