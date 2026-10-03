package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

var sparepartsData []Sparepart
var stokData []DataStok
var hargaData []DataHarga
var lokasiData []LokasiGudang
var supplierData []DataSupplier

func main() {

	// ==============================
	// MEMBACA SEMUA DATA JSON
	// ==============================

	spareparts, stok, harga, lokasi, suppliers, err := readAllData()

	if err != nil {
		log.Fatal(err)
	}

	sparepartsData = spareparts
	stokData = stok
	hargaData = harga
	lokasiData = lokasi
	supplierData = suppliers

	go startAPI()

	// Membuat reader untuk membaca input keyboard
	reader := bufio.NewReader(os.Stdin)

	// ==============================
	// MENU UTAMA
	// ==============================

	for {

		fmt.Println()
		fmt.Println("======================================")
		fmt.Println("       INVENTORY SPAREPART MOTOR")
		fmt.Println("======================================")
		fmt.Println("1. Tampilkan Semua Barang")
		fmt.Println("2. Tambah Barang")
		fmt.Println("3. Cari Barang")
		fmt.Println("4. Login")
		fmt.Println("0. Keluar")
		fmt.Println("======================================")

		pilihan := readLine(reader, "Pilih Menu: ")

		// ==============================
		// MEMPROSES PILIHAN MENU
		// ==============================

		switch pilihan {

		case "1":

			// Menampilkan semua data
			showInventory(
				spareparts,
				stok,
				harga,
				lokasi,
				suppliers,
			)

		case "2":

			// Menambahkan barang
			spareparts, stok, harga, lokasi, suppliers =
				tambahBarang(
					reader,
					spareparts,
					stok,
					harga,
					lokasi,
					suppliers,
				)

		case "3":

			// Menu pencarian barang
			cariBarang(reader, spareparts, stok, harga, lokasi, suppliers)

		case "4":

			// Menu login
			credentialCheck()

		case "0":

			fmt.Println()
			fmt.Println("Program selesai.")
			return

		default:

			fmt.Println()
			fmt.Println("Menu tidak tersedia.")
		}
	}
}
