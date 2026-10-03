package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// ============================================================
// LOGIN
// ============================================================

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthUser struct {
	ID          string      `json:"id"`
	Credentials Credentials `json:"credentials"`
}

var users = []AuthUser{
	{
		ID: "USR-001",
		Credentials: Credentials{
			Username: "admin",
			Password: "admin123",
		},
	},
	{
		ID: "USR-002",
		Credentials: Credentials{
			Username: "kasir01",
			Password: "kasir123",
		},
	},
	{
		ID: "USR-003",
		Credentials: Credentials{
			Username: "user01",
			Password: "user123",
		},
	},
}

func login(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var request Credentials

	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		http.Error(
			w,
			"JSON tidak valid",
			http.StatusBadRequest,
		)
		return
	}

	for _, user := range users {

		if user.Credentials.Username == request.Username &&
			user.Credentials.Password == request.Password {

			json.NewEncoder(w).Encode(map[string]interface{}{
				"success":  true,
				"message":  "Login berhasil",
				"user_id":  user.ID,
				"username": user.Credentials.Username,
			})

			return
		}
	}

	w.WriteHeader(http.StatusUnauthorized)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": false,
		"message": "Username atau password salah",
	})
}

// ============================================================
// GET - MENAMPILKAN SEMUA SPAREPART
// ============================================================

func getSpareparts(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(sparepartsData)
}

// ============================================================
// POST - TAMBAH SPAREPART
// ============================================================

func postSparepart(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var request TambahBarangRequest

	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		http.Error(w, "JSON tidak valid", http.StatusBadRequest)
		return
	}

	// Membuat ID baru
	id := nextID(sparepartsData)

	kodeBarang := fmt.Sprintf("BRG-%03d", id)

	// Membuat data sparepart
	dataSparepart := Sparepart{
		No:          id,
		KodeBarang:  kodeBarang,
		Identitas:   request.Identitas,
		Spesifikasi: request.Spesifikasi,
	}

	// Membuat data stok
	dataStok := DataStok{
		KodeBarang: kodeBarang,
		Stok:       request.Stok,
	}

	// Membuat data harga
	dataHarga := DataHarga{
		KodeBarang: kodeBarang,
		Harga:      request.Harga,
	}

	// Membuat data lokasi
	dataLokasi := LokasiGudang{
		KodeBarang:   kodeBarang,
		LokasiGudang: request.Lokasi,
	}

	// Membuat data supplier
	dataSupplier := DataSupplier{
		KodeBarang: kodeBarang,
		Supplier:   request.Supplier,
	}

	// Menambahkan ke memory
	sparepartsData = append(sparepartsData, dataSparepart)
	stokData = append(stokData, dataStok)
	hargaData = append(hargaData, dataHarga)
	lokasiData = append(lokasiData, dataLokasi)
	supplierData = append(supplierData, dataSupplier)

	// Menyimpan ke JSON
	if err := simpanSemuaData(); err != nil {
		http.Error(
			w,
			"Gagal menyimpan data: "+err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":     "Barang berhasil ditambahkan",
		"kode_barang": kodeBarang,
		"data":        dataSparepart,
	})
}

// ============================================================
// PUT - EDIT SPAREPART
// ============================================================

func putSparepart(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	// Mengambil kode dari URL
	kodeBarang := strings.TrimPrefix(
		r.URL.Path,
		"/api/sparepart/",
	)

	if kodeBarang == "" {
		http.Error(
			w,
			"Kode barang tidak ditemukan",
			http.StatusBadRequest,
		)
		return
	}

	// Mencari barang
	index := -1

	for i, item := range sparepartsData {

		if strings.EqualFold(item.KodeBarang, kodeBarang) {
			index = i
			break
		}
	}

	if index == -1 {

		http.Error(
			w,
			"Barang tidak ditemukan",
			http.StatusNotFound,
		)

		return
	}

	// Membaca JSON
	var request TambahBarangRequest

	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {

		http.Error(
			w,
			"JSON tidak valid",
			http.StatusBadRequest,
		)

		return
	}

	// Nomor barang tetap
	no := sparepartsData[index].No

	// Update sparepart
	sparepartsData[index] = Sparepart{
		No:          no,
		KodeBarang:  sparepartsData[index].KodeBarang,
		Identitas:   request.Identitas,
		Spesifikasi: request.Spesifikasi,
	}

	// Update stok
	for i := range stokData {

		if strings.EqualFold(
			stokData[i].KodeBarang,
			kodeBarang,
		) {

			stokData[i].Stok = request.Stok

			break
		}
	}

	// Update harga
	for i := range hargaData {

		if strings.EqualFold(
			hargaData[i].KodeBarang,
			kodeBarang,
		) {

			hargaData[i].Harga = request.Harga

			break
		}
	}

	// Update lokasi
	for i := range lokasiData {

		if strings.EqualFold(
			lokasiData[i].KodeBarang,
			kodeBarang,
		) {

			lokasiData[i].LokasiGudang = request.Lokasi

			break
		}
	}

	// Update supplier
	for i := range supplierData {

		if strings.EqualFold(
			supplierData[i].KodeBarang,
			kodeBarang,
		) {

			supplierData[i].Supplier = request.Supplier

			break
		}
	}

	// Simpan semua JSON
	if err := simpanSemuaData(); err != nil {

		http.Error(
			w,
			"Gagal menyimpan data: "+err.Error(),
			http.StatusInternalServerError,
		)

		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Barang berhasil diperbarui",
		"data":    sparepartsData[index],
	})
}

// ============================================================
// DELETE - HAPUS SPAREPART
// ============================================================

func deleteSparepart(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	kodeBarang := strings.TrimPrefix(
		r.URL.Path,
		"/api/sparepart/",
	)

	if kodeBarang == "" {

		http.Error(
			w,
			"Kode barang tidak ditemukan",
			http.StatusBadRequest,
		)

		return
	}

	// Cari index
	index := -1

	for i, item := range sparepartsData {

		if strings.EqualFold(
			item.KodeBarang,
			kodeBarang,
		) {

			index = i
			break
		}
	}

	if index == -1 {

		http.Error(
			w,
			"Barang tidak ditemukan",
			http.StatusNotFound,
		)

		return
	}

	// Hapus dari sparepart
	sparepartsData = append(
		sparepartsData[:index],
		sparepartsData[index+1:]...,
	)

	// Hapus stok
	for i, item := range stokData {

		if strings.EqualFold(
			item.KodeBarang,
			kodeBarang,
		) {

			stokData = append(
				stokData[:i],
				stokData[i+1:]...,
			)

			break
		}
	}

	// Hapus harga
	for i, item := range hargaData {

		if strings.EqualFold(
			item.KodeBarang,
			kodeBarang,
		) {

			hargaData = append(
				hargaData[:i],
				hargaData[i+1:]...,
			)

			break
		}
	}

	// Hapus lokasi
	for i, item := range lokasiData {

		if strings.EqualFold(
			item.KodeBarang,
			kodeBarang,
		) {

			lokasiData = append(
				lokasiData[:i],
				lokasiData[i+1:]...,
			)

			break
		}
	}

	// Hapus supplier
	for i, item := range supplierData {

		if strings.EqualFold(
			item.KodeBarang,
			kodeBarang,
		) {

			supplierData = append(
				supplierData[:i],
				supplierData[i+1:]...,
			)

			break
		}
	}

	// Simpan
	if err := simpanSemuaData(); err != nil {

		http.Error(
			w,
			"Gagal menyimpan data: "+err.Error(),
			http.StatusInternalServerError,
		)

		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":     "Barang berhasil dihapus",
		"kode_barang": kodeBarang,
	})
}

// ============================================================
// SIMPAN SEMUA DATA
// ============================================================

func simpanSemuaData() error {

	if err := saveJSON(
		sparepartFile,
		sparepartsData,
	); err != nil {
		return err
	}

	if err := saveJSON(
		stokFile,
		stokData,
	); err != nil {
		return err
	}

	if err := saveJSON(
		hargaFile,
		hargaData,
	); err != nil {
		return err
	}

	if err := saveJSON(
		lokasiFile,
		lokasiData,
	); err != nil {
		return err
	}

	if err := saveJSON(
		supplierFile,
		supplierData,
	); err != nil {
		return err
	}

	return nil
}

// ============================================================
// START API
// ============================================================

func startAPI() {

	mux := http.NewServeMux()

	// ============================================================
	// LOGIN
	// ============================================================

	mux.HandleFunc("/api/login", func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			http.Error(
				w,
				"Method tidak didukung",
				http.StatusMethodNotAllowed,
			)
			return
		}

		login(w, r)
	})

	// ============================================================
	// GET SEMUA SPAREPART
	// ============================================================

	mux.HandleFunc("/api/sparepart", func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodGet {
			http.Error(
				w,
				"Method tidak didukung",
				http.StatusMethodNotAllowed,
			)
			return
		}

		getSpareparts(w, r)
	})

	// ============================================================
	// POST TAMBAH SPAREPART
	// ============================================================

	mux.HandleFunc("/api/sparepart/tambah", func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			http.Error(
				w,
				"Method tidak didukung",
				http.StatusMethodNotAllowed,
			)
			return
		}

		postSparepart(w, r)
	})

	// ============================================================
	// PUT DAN DELETE
	// ============================================================

	mux.HandleFunc("/api/sparepart/", func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {

		case http.MethodPut:
			putSparepart(w, r)

		case http.MethodDelete:
			deleteSparepart(w, r)

		default:
			http.Error(
				w,
				"Method tidak didukung",
				http.StatusMethodNotAllowed,
			)
		}
	})

	// ============================================================
	// WEBSITE
	// ============================================================

	mux.Handle(
		"/",
		http.FileServer(
			http.Dir("./web"),
		),
	)

	fmt.Println("======================================")
	fmt.Println("       API SERVER BERJALAN")
	fmt.Println("======================================")
	fmt.Println("WEB    : http://localhost:8080")
	fmt.Println("GET    : http://localhost:8080/api/sparepart")
	fmt.Println("POST   : http://localhost:8080/api/sparepart/tambah")
	fmt.Println("PUT    : http://localhost:8080/api/sparepart/{kode}")
	fmt.Println("DELETE : http://localhost:8080/api/sparepart/{kode}")
	fmt.Println("======================================")

	err := http.ListenAndServe(":8080", mux)

	if err != nil {
		fmt.Println("Server error:", err)
	}
}
