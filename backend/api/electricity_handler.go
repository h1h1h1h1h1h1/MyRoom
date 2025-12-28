package api

import (
	"database/sql"
	"net/http"
	"strconv"

	"electricity-management-backend/database"
	"electricity-management-backend/models"

	"github.com/gin-gonic/gin"
)

type ElectricityHandler struct{}

func NewElectricityHandler() *ElectricityHandler {
	return &ElectricityHandler{}
}

func (h *ElectricityHandler) GetUsage(c *gin.Context) {
	customerIDStr := c.Query("customer_id")
	customerID, err := strconv.Atoi(customerIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customer_id"})
		return
	}

	rows, err := database.DB.Query("SELECT id, customer_number_id, reading_date, current_reading, previous_reading, consumption, reading_type, created_at FROM meter_readings WHERE customer_number_id = ? ORDER BY reading_date DESC", customerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	var readings []models.MeterReading
	for rows.Next() {
		var r models.MeterReading
		var prevReading sql.NullFloat64
		if err := rows.Scan(&r.ID, &r.CustomerNumberID, &r.ReadingDate, &r.CurrentReading, &prevReading, &r.Consumption, &r.ReadingType, &r.CreatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan reading"})
			return
		}
		if prevReading.Valid {
			r.PreviousReading = prevReading.Float64
		}
		readings = append(readings, r)
	}

	if readings == nil {
		readings = []models.MeterReading{}
	}

	c.JSON(http.StatusOK, readings)
}

func (h *ElectricityHandler) GetPayments(c *gin.Context) {
	customerIDStr := c.Query("customer_id")
	customerID, err := strconv.Atoi(customerIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid customer_id"})
		return
	}

	rows, err := database.DB.Query("SELECT id, customer_number_id, billing_month, consumption, unit_price, amount, due_date, paid_amount, payment_status, paid_date, created_at FROM electricity_bills WHERE customer_number_id = ? ORDER BY billing_month DESC", customerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	defer rows.Close()

	var bills []models.ElectricityBill
	for rows.Next() {
		var bill models.ElectricityBill
		var paidDate sql.NullTime
		if err := rows.Scan(&bill.ID, &bill.CustomerNumberID, &bill.BillingMonth, &bill.Consumption, &bill.UnitPrice, &bill.Amount, &bill.DueDate, &bill.PaidAmount, &bill.PaymentStatus, &paidDate, &bill.CreatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan bill"})
			return
		}
		if paidDate.Valid {
			bill.PaidDate = paidDate.Time
		}
		bills = append(bills, bill)
	}

	if bills == nil {
		bills = []models.ElectricityBill{}
	}

	c.JSON(http.StatusOK, bills)
}

func (h *ElectricityHandler) Pay(c *gin.Context) {
	var req struct {
		BillID int     `json:"bill_id"`
		Amount float64 `json:"amount"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update bill status
	_, err := database.DB.Exec("UPDATE electricity_bills SET paid_amount = paid_amount + ?, payment_status = 'paid', paid_date = NOW() WHERE id = ?", req.Amount, req.BillID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process payment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment successful"})
}
