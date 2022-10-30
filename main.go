package main

import (
	"context"

	"bootcamp-day-12/connection"

	"bootcamp-day-12/middleware"

	"fmt"

	"html/template"

	"log"

	"net/http"

	"time"

	"strings"

	"strconv"

	"github.com/gorilla/mux"

	"golang.org/x/crypto/bcrypt"

	"github.com/gorilla/sessions"
)

type Session struct {
	Title     string
	IsLogin   bool
	UserName  string
	UserId    int
	FlashData string
} // Sebagai struktur untuk menampung data sesi login

var Data = Session{} // Untuk menghubungkan data dengan sesi login

type Project struct {
	Id           int
	Title        string
	Start_date   time.Time
	Format_start string
	End_date     time.Time
	Format_end   string
	Description  string
	Technologies []string
	Image        string
	User_id      int
	Duration     string
	Author       string
	IsLogin      bool
} // Sebagai struktur untuk menampung data terkait project

type User struct {
	Id       int
	Name     string
	Email    string
	Password string
} // Sebagai struktur untuk menampung data terkait user

func main() {
	route := mux.NewRouter() // sebuah package yang di impor dari gorilla mux untuk fitur routing

	connection.DatabaseConnect() // Untuk melakukan koneksi dengan database

	// Untuk mengakses direktori dimana kita menyimpan dokumen yang dibutuhkan
	route.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))
	route.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))

	// untuk melakukan navigasi terhadap masing-masing page
	route.HandleFunc("/", home).Methods("GET")
	route.HandleFunc("/contact", contact).Methods("GET")

	route.HandleFunc("/formProject", formProject).Methods("GET")
	route.HandleFunc("/projectDetail/{id}", projectDetail).Methods("GET")
	route.HandleFunc("/addProject", middleware.UploadFile(addProject)).Methods("POST") // untuk nenjembatani fitur upload file ke local storage

	route.HandleFunc("/deleteProject/{id}", deleteProject).Methods("GET")
	route.HandleFunc("/updateForm/{id}", updateForm).Methods("GET")
	route.HandleFunc("/updateProject/{id}", middleware.UploadFile(updateProject)).Methods("POST") // untuk nenjembatani fitur upload file ke local

	route.HandleFunc("/formRegister", formRegister).Methods("GET")
	route.HandleFunc("/register", register).Methods("POST")

	route.HandleFunc("/formLogin", formLogin).Methods("GET")
	route.HandleFunc("/login", login).Methods("POST")

	route.HandleFunc("/logout", logout).Methods("GET")

	fmt.Println("Server Running on port 5000")
	http.ListenAndServe("localhost:5000", route) // To run the code on local server
}

// ---------------------------------------------------------------------
// BARIS KODE UNTUK MELAKUKAN PEMANGGILAN RENDER PAGE
// ---------------------------------------------------------------------

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmpl, err := template.ParseFiles("views/index.html")
	// ERROR HANDLING RENDER HTML TEMPLATE
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	} else {
		// Untuk mendapatkan data sesi login dari cookies browser
		store := sessions.NewCookieStore([]byte("SESSION_KEY"))
		session, _ := store.Get(r, "SESSION_KEY")

		// Untuk melakukan pengecekan terhadap status login dan menampilkan alert
		if session.Values["IsLogin"] != true {
			Data.IsLogin = false

			fm := session.Flashes("message")

			var flashes []string

			// Melakukan pengecekan apabila terdapat alert
			if len(fm) > 0 {
				session.Save(r, w)
				for _, f1 := range fm {
					flashes = append(flashes, f1.(string))
				}
			}
			// Menambahkan data ke sesi login
			Data.FlashData = strings.Join(flashes, "")
			Data.UserId = 0
		} else {

			fm := session.Flashes("message")

			var flashes []string

			// Melakukan pengecekan apabila terdapat alert
			if len(fm) > 0 {
				session.Save(r, w)
				for _, f1 := range fm {
					flashes = append(flashes, f1.(string))
				}
			}

			// Menambahkan data ke dalam sesi login
			Data.FlashData = strings.Join(flashes, "")
			Data.IsLogin = session.Values["IsLogin"].(bool)
			Data.UserName = session.Values["UserName"].(string)
			Data.UserId = session.Values["UserId"].(int)
		}

		var projectRender []Project // variable untuk memanggil slice project
		each := Project{}           // variable untuk mengakses struct project

		if !Data.IsLogin {
			// Membuat koneksi dengan antara tabel project dan tabel user
			rows, _ := connection.Conn.Query(context.Background(), "SELECT tb_projects.id, title, start_date, end_date, description, technologies, image, user_id, name FROM public.tb_projects LEFT JOIN public.tb_users ON public.tb_projects.user_id = public.tb_users.id")

			// Melakukan scanning terhadap value yang akan ditampilkan
			for rows.Next() {
				err := rows.Scan(&each.Id, &each.Title, &each.Start_date, &each.End_date, &each.Description, &each.Technologies, &each.Image, &each.User_id, &each.Author)

				if err != nil {
					fmt.Println(err.Error())
					return
				} else {
					// Melakukan pengisian value database terhadap struct Project
					each := Project{
						Id:           each.Id,
						Title:        each.Title,
						Duration:     DurationCount(each.Start_date, each.End_date),
						Description:  each.Description,
						Technologies: each.Technologies,
						Image:        each.Image,
						User_id:      each.User_id,
						Author:       each.Author,
					}
					projectRender = append(projectRender, each) // Memasukkan nilai variable each ke dalam variable projectRender yang merupakan []Project
				}
			}
			// Nilai yang sudah dideklarasikan sebelumnya ditampung dalam variable respData sebelum dieksekusi oleh tmpl.Execute
			respData := map[string]interface{}{
				"projectRender": projectRender,
				"Data":          Data,
			}
			w.WriteHeader(http.StatusOK)
			tmpl.Execute(w, respData)

		} else {

			// Membuat koneksi dengan antara tabel project dan tabel user dan hanya mengakses user_id yang sedang login
			rows, _ := connection.Conn.Query(context.Background(), "SELECT tb_projects.id, title, start_date, end_date, description, technologies, image, user_id, name FROM public.tb_projects LEFT JOIN public.tb_users ON public.tb_projects.user_id = public.tb_users.id WHERE user_id = $1", Data.UserId)

			// Melakukan scanning terhadap value yang akan ditampilkan
			for rows.Next() {

				err := rows.Scan(&each.Id, &each.Title, &each.Start_date, &each.End_date, &each.Description, &each.Technologies, &each.Image, &each.User_id, &each.Author)

				if err != nil {
					fmt.Println(err.Error())
					return
				} else {
					// Melakukan pengisian value database terhadap struct Project
					each := Project{
						Id:           each.Id,
						Title:        each.Title,
						Duration:     DurationCount(each.Start_date, each.End_date),
						Description:  each.Description,
						Technologies: each.Technologies,
						Image:        each.Image,
						User_id:      each.User_id,
						Author:       each.Author,
					}
					projectRender = append(projectRender, each) // Memasukkan nilai variable each ke dalam variable projectRender yang merupakan []Project
				}
			}
			// Nilai yang sudah dideklarasikan sebelumnya ditampung dalam variable respData sebelum dieksekusi oleh tmpl.Execute
			respData := map[string]interface{}{
				"projectRender": projectRender,
				"Data":          Data,
			}
			w.WriteHeader(http.StatusOK)
			tmpl.Execute(w, respData)
		}

	}
}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmpl, err := template.ParseFiles("views/contact.html")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	} else {
		// Untuk mendapatkan data sesi login dari cookies browser
		store := sessions.NewCookieStore([]byte("SESSION_KEY"))
		session, _ := store.Get(r, "SESSION_KEY")
		// Melakukan pengecekan status login
		if session.Values["IsLogin"] != true {
			Data.IsLogin = false
		} else {
			Data.IsLogin = session.Values["IsLogin"].(bool)
			Data.UserName = session.Values["UserName"].(string)
			Data.UserId = session.Values["UserId"].(int)
		}
		respData := map[string]interface{}{
			"Data": Data,
		}

		w.WriteHeader(http.StatusOK)
		tmpl.Execute(w, respData)
	}
}

func formProject(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmpl, err := template.ParseFiles("views/form-project.html")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	} else {
		// Mengakses cookies browser
		store := sessions.NewCookieStore([]byte("SESSION_KEY"))
		session, _ := store.Get(r, "SESSION_KEY")
		// mengecek status login
		if session.Values["IsLogin"] != true {
			Data.IsLogin = false
		} else {
			Data.IsLogin = session.Values["IsLogin"].(bool)
			Data.UserName = session.Values["UserName"].(string)
			Data.UserId = session.Values["UserId"].(int)
		}
		respData := map[string]interface{}{
			"Data": Data,
		}
		w.WriteHeader(http.StatusOK)
		tmpl.Execute(w, respData)
	}
}

func formRegister(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmpl, err := template.ParseFiles("views/form-register.html")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	} else {
		// Mengakses cookies browser
		store := sessions.NewCookieStore([]byte("SESSION_KEY"))
		session, _ := store.Get(r, "SESSION_KEY")
		// Mengecek login status
		if session.Values["IsLogin"] != true {
			Data.IsLogin = false

			fm := session.Flashes("message")

			var flashes []string
			// Mengecek apakah terdapat alert
			if len(fm) > 0 {
				session.Save(r, w)
				for _, f1 := range fm {
					flashes = append(flashes, f1.(string))
				}
			}
			// Memasukkan data alert ke sesi login
			Data.FlashData = strings.Join(flashes, "")

			respData := map[string]interface{}{
				"Data": Data,
			}

			w.WriteHeader(http.StatusOK)
			tmpl.Execute(w, respData)

		} else {
			Data.IsLogin = session.Values["IsLogin"].(bool)
			Data.UserName = session.Values["UserName"].(string)
			Data.UserId = session.Values["UserId"].(int)

			respData := map[string]interface{}{
				"Data": Data,
			}

			w.WriteHeader(http.StatusOK)
			tmpl.Execute(w, respData)
		}
	}
}

func formLogin(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmpl, err := template.ParseFiles("views/form-login.html")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	} else {
		// Mengakses cookies browser
		store := sessions.NewCookieStore([]byte("SESSION_KEY"))
		session, _ := store.Get(r, "SESSION_KEY")
		// Mengecek status sesi login
		if session.Values["IsLogin"] != true {
			Data.IsLogin = false

			fm := session.Flashes("message")

			var flashes []string
			// Mengecek apakah terdapat alert
			if len(fm) > 0 {
				session.Save(r, w)
				for _, f1 := range fm {
					flashes = append(flashes, f1.(string))
				}
			}
			// Menambahkan data alert ke dalam sesi login
			Data.FlashData = strings.Join(flashes, "")

			respData := map[string]interface{}{
				"Data": Data,
			}
			tmpl.Execute(w, respData)

		} else {
			Data.IsLogin = session.Values["IsLogin"].(bool)
			Data.UserName = session.Values["UserName"].(string)
			Data.UserId = session.Values["UserId"].(int)

			respData := map[string]interface{}{
				"Data": Data,
			}

			w.WriteHeader(http.StatusOK)
			tmpl.Execute(w, respData)
		}
	}
}

func projectDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmpl, err := template.ParseFiles("views/project-page.html")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	} else {
		id, _ := strconv.Atoi(mux.Vars(r)["id"]) // mengkonversi tipe data id dari string ke int

		projectDetail := Project{} // variable untuk menakses struct Project

		// Membuat koneksi seleksi dengan database dan melakukan scanning terhadap nilainya
		err = connection.Conn.QueryRow(context.Background(), `SELECT *
		FROM public.tb_projects WHERE "id" = $1`, id).Scan(&projectDetail.Id, &projectDetail.Title, &projectDetail.Start_date, &projectDetail.End_date, &projectDetail.Description, &projectDetail.Technologies, &projectDetail.Image, &projectDetail.User_id)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("message : " + err.Error()))

		} else {
			// Mengakses cookies browser
			store := sessions.NewCookieStore([]byte("SESSION_KEY"))
			session, _ := store.Get(r, "SESSION_KEY")

			// Mengecek status sesi login
			if session.Values["IsLogin"] != true {
				Data.IsLogin = false
			} else {
				Data.IsLogin = session.Values["IsLogin"].(bool)
				Data.UserName = session.Values["UserName"].(string)
				Data.UserId = session.Values["UserId"].(int)
			}

			projectDetail.Format_start = projectDetail.Start_date.Format("2 January 2006") // mengubah format penuliasan tanggal
			projectDetail.Format_end = projectDetail.End_date.Format("2 January 2006")

			// Memanggil nilai yang ada pada project struct untuk nantinya ditampilkan di project page
			projectDetail := Project{
				Id:           projectDetail.Id,
				Title:        projectDetail.Title,
				Format_start: projectDetail.Format_start,
				Format_end:   projectDetail.Format_end,
				Duration:     DurationCount(projectDetail.Start_date, projectDetail.End_date),
				Description:  projectDetail.Description,
				Technologies: projectDetail.Technologies,
				Image:        projectDetail.Image,
				User_id:      projectDetail.User_id,
			}
			respData := map[string]interface{}{
				"projectDetail": projectDetail,
				"Data":          Data,
			}
			w.WriteHeader(http.StatusOK)
			tmpl.Execute(w, respData)
		}
	}
}

func updateForm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmpl, err := template.ParseFiles("views/form-update.html")

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
		return
	} else {
		id, _ := strconv.Atoi(mux.Vars(r)["id"]) // Melakukan konversi tipe data id dari string ke int
		updateProject := Project{}

		// Melakukan koneksi dengan database dan melakukan scanning terhadap nilainya
		err = connection.Conn.QueryRow(context.Background(), `SELECT * FROM public.tb_projects WHERE "id" = $1`, id).Scan(&updateProject.Id, &updateProject.Title, &updateProject.Start_date, &updateProject.End_date, &updateProject.Description, &updateProject.Technologies, &updateProject.Image, &updateProject.User_id)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("message : " + err.Error()))
			return
		}
		// Mengakses cookies browser
		store := sessions.NewCookieStore([]byte("SESSION_KEY"))
		session, _ := store.Get(r, "SESSION_KEY")
		// Mengecek status login
		if session.Values["IsLogin"] != true {
			Data.IsLogin = false
		} else {
			Data.IsLogin = session.Values["IsLogin"].(bool)
			Data.UserName = session.Values["UserName"].(string)
			Data.UserId = session.Values["UserId"].(int)
		}
		// mengakses project struct yang akan diupdate
		updateProject = Project{
			Id:           updateProject.Id,
			Title:        updateProject.Title,
			Start_date:   updateProject.Start_date,
			End_date:     updateProject.End_date,
			Duration:     DurationCount(updateProject.Start_date, updateProject.End_date),
			Description:  updateProject.Description,
			Technologies: updateProject.Technologies,
			Image:        updateProject.Image,
			User_id:      updateProject.User_id,
		}

		respData := map[string]interface{}{
			"updateProject": updateProject,
			"Data":          Data,
		}

		w.WriteHeader(http.StatusOK)
		tmpl.Execute(w, respData)
	}
}

// --------------------------------------------------------------------
// BARIS KODE UNTUK AUTENTIFIKASI DATA USER
// --------------------------------------------------------------------

func register(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm() // Untuk melakukan parsing data terhadap form pada view html

	if err != nil {
		log.Fatal(err)
	} else {
		Name := r.PostForm.Get("inputName")
		Email := r.PostForm.Get("inputEmail")
		Password := r.PostForm.Get("inputPassword")

		// Melakukan enkripsi password dengan package bcrypt
		PasswordHash, _ := bcrypt.GenerateFromPassword([]byte(Password), 10)

		// melakukan koneksi dengan database dan memasukkan data ke table users
		_, err = connection.Conn.Exec(context.Background(), `INSERT INTO public.tb_users("name", "email", "password") VALUES($1, $2, $3)`, Name, Email, PasswordHash)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("message : " + err.Error()))
			return
		} else {
			// Mengakses data cookies browser
			var store = sessions.NewCookieStore([]byte("SESSION_KEY"))
			session, _ := store.Get(r, "SESSION_KEY")

			// menambahkan data dan alert terhadap halaman login
			session.Values["IsLogin"] = false
			session.AddFlash("Akun berhasil didaftarkan, silahkan login!", "message")
			session.Save(r, w)

			http.Redirect(w, r, "/formLogin", http.StatusMovedPermanently)
		}
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	// Mengakses cookies browser
	store := sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")

	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	} else {
		Email := r.PostForm.Get("inputEmail")
		Password := r.PostForm.Get("inputPassword")

		// Mengakses struck user yang login
		userLogin := User{}
		// Melakukan koneksi dengan database dan melakukan scanning terhadap nilainya
		err = connection.Conn.QueryRow(context.Background(), `SELECT * FROM public.tb_users WHERE "email" = $1`,
			Email).Scan(&userLogin.Id, &userLogin.Name, &userLogin.Email, &userLogin.Password)

		if err != nil {
			var store = sessions.NewCookieStore([]byte("SESSION_KEY"))
			session, _ := store.Get(r, "SESSION_KEY")

			// memanggil alert Jika email yang dimasukkan salah
			session.Values["IsLogin"] = false
			session.AddFlash("Mohon maaf, email Anda tidak terdaftar!", "message")
			session.Save(r, w)

			http.Redirect(w, r, "/formLogin", http.StatusMovedPermanently)
			return
		} else {
			// Melakukan pengecekan terhadap nilai password
			err = bcrypt.CompareHashAndPassword([]byte(userLogin.Password), []byte(Password))
			// Memanggil alert jika password tidak sesuai database
			if err != nil {
				var store = sessions.NewCookieStore([]byte("SESSION_KEY"))
				session, _ := store.Get(r, "SESSION_KEY")

				session.Values["IsLogin"] = false
				session.AddFlash("Mohon maaf, password Anda salah!", "message")
				session.Save(r, w)

				http.Redirect(w, r, "/formLogin", http.StatusMovedPermanently)
				return
			} else {
				// Memasukkan data ke sesi login jika email dan password sudah sesuai database
				session.Values["IsLogin"] = true
				session.Values["UserName"] = userLogin.Name
				session.Values["UserId"] = userLogin.Id

				session.Options.MaxAge = 10800 // 10800 Detik : 3600 detik = Maka durasi login 3 jam

				// Manambahkan alert ke dalam sesi login
				session.AddFlash("Berhasil Login! Selamat Datang "+userLogin.Name, "message")
				session.Save(r, w)

				http.Redirect(w, r, "/", http.StatusMovedPermanently)
			}
		}
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	// Mengakses cookies browser
	var store = sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")

	session.AddFlash("Anda telah Logout!", "message") // menambahkan alert pada sesi login

	session.Options.MaxAge = -1 // menghapus durasi login

	session.Values["IsLogin"] = false
	session.Values["UserId"] = 0

	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// -------------------------------------------------------------------
// BARIS KODE UNTUK MODIFIKASI TERKAIT DATA PROJECT
// -------------------------------------------------------------------

func addProject(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		log.Fatal(err)
	} else {
		// Mengakses cookies browser
		store := sessions.NewCookieStore([]byte("SESSION_KEY"))
		session, _ := store.Get(r, "SESSION_KEY")
		// mengecek status login
		if session.Values["IsLogin"] != true {
			Data.IsLogin = false
		} else {
			Data.IsLogin = session.Values["IsLogin"].(bool)
			Data.UserName = session.Values["UserName"].(string)
			Data.UserId = session.Values["UserId"].(int)
		}

		// Mengambil nilai dari form
		projectTitle := r.PostForm.Get("input-title")
		projectStart := r.PostForm.Get("start-date")
		projectEnd := r.PostForm.Get("end-date")
		projectDescription := r.PostForm.Get("project-description")
		projectTechnologies := []string{r.PostForm.Get("reactjs"), r.PostForm.Get("nodejs"), r.PostForm.Get("javascript"), r.PostForm.Get("golang")}

		dataContext := r.Context().Value("dataFile") // data context sesuai dengan function middleware
		projectImage := dataContext.(string)

		// Mengakses UserId
		UserId := Data.UserId

		// Melakukan koneksi dengan database dan memasukkan nilai yang telah diambil dari form ke database
		_, err = connection.Conn.Exec(context.Background(), `INSERT INTO public.tb_projects( "title", "start_date", "end_date", "description", "technologies", "image", "user_id") VALUES ( $1, $2, $3, $4, $5, $6, $7)`, projectTitle, projectStart, projectEnd, projectDescription, projectTechnologies, projectImage, UserId)

		// Memanggil alert untuk ditambahkan ke sesi login
		if err != nil {
			session.AddFlash("Project tidak dapat ditambahkan!", "message")
			session.Save(r, w)

			http.Redirect(w, r, "/", http.StatusMovedPermanently)
		} else {
			session.AddFlash("Project berhasil ditambahkan!", "message")
			session.Save(r, w)

			http.Redirect(w, r, "/", http.StatusMovedPermanently)
		}
	}
}

func updateProject(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm() // Melakukan parsing terhadap data form pada views html

	if err != nil {
		log.Fatal(err)
	} else {
		// Mengakses cookies browser
		store := sessions.NewCookieStore([]byte("SESSION_KEY"))
		session, _ := store.Get(r, "SESSION_KEY")
		// Mengecek status sesi login
		if session.Values["IsLogin"] != true {
			Data.IsLogin = false
		} else {
			Data.IsLogin = session.Values["IsLogin"].(bool)
			Data.UserName = session.Values["UserName"].(string)
			Data.UserId = session.Values["UserId"].(int)
		}

		id, _ := strconv.Atoi(mux.Vars(r)["id"]) // melakukan konversi tipe data id dari string ke int

		// mengambil nilai dari input form
		projectTitle := r.PostForm.Get("input-title")
		projectStart := r.PostForm.Get("start-date")
		projectEnd := r.PostForm.Get("end-date")
		projectDescription := r.PostForm.Get("project-description")
		projectTechnologies := []string{r.PostForm.Get("reactjs"), r.PostForm.Get("nodejs"), r.PostForm.Get("javascript"), r.PostForm.Get("golang")}

		dataContext := r.Context().Value("dataFile") // data context sesuai dengan middleware
		projectImage := dataContext.(string)

		// Melakukan koneksi dengan database dan juga melakukan update dari data yang telah di input form
		_, err = connection.Conn.Exec(context.Background(), `UPDATE public.tb_projects SET "title"=$1, "start_date"=$2, "end_date"=$3, "description"=$4, "technologies"=$5, "image"=$6 WHERE "id"=$7`, projectTitle, projectStart, projectEnd, projectDescription, projectTechnologies, projectImage, id)

		// Memanggil alert terhadap sesi login
		if err != nil {
			session.AddFlash("Update project tidak berhasil!", "message")
			session.Save(r, w)

			http.Redirect(w, r, "/", http.StatusMovedPermanently)
		} else {
			session.AddFlash("Update project berhasil!", "message")
			session.Save(r, w)

			http.Redirect(w, r, "/", http.StatusMovedPermanently)
		}
	}
}

func deleteProject(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"]) // melakukan konversi terhadap tipe data id dari string ke int

	// mengakses cookies browser
	store := sessions.NewCookieStore([]byte("SESSION_KEY"))
	session, _ := store.Get(r, "SESSION_KEY")

	// Melakukan koneksi terhadap database dan menghapus data berdasarkan id yang diseleksi
	_, err := connection.Conn.Exec(context.Background(), `DELETE FROM public.tb_projects WHERE "id" = $1`, id)

	// Memanggil alert terhadap sesi login
	if err != nil {
		session.AddFlash("Tidak dapat menghapus project!", "message")
		session.Save(r, w)

		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	} else {
		session.AddFlash("Project berhasil terhapus!", "message")
		session.Save(r, w)

		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	}
}

func DurationCount(Start_date time.Time, End_date time.Time) string {
	days := End_date.Sub(Start_date).Hours() / 24 // Selisih tanggal, output diubah ke jam, jam berubah jadi hari
	var duration string

	if days >= 30 {
		if (days / 30) == 1 {
			duration = "duration 1 month"
		} else {
			duration = "duration " + strconv.Itoa(int(days/30)) + " months" // mengubah data ke string untuk hasil bulan
		}
	} else {
		if days <= 1 {
			duration = "1 day"
		} else {
			duration = "duration " + strconv.Itoa(int(days)) + " days" // mengubah data ke string untuk hasil hari
		}
	}
	return duration
}
