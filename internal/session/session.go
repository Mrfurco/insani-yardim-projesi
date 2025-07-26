package session

import (
	"encoding/gob" // GOB KÜTÜPHANESİNİ İMPORT EDİYORUZ

	"github.com/MrFurco/insani-yardim-projesi/internal/model" // MODEL PAKETİMİZİ İMPORT EDİYORUZ
	"github.com/gorilla/sessions"
)

// Store, tüm uygulama tarafından kullanılacak olan tek ve merkezi oturum yöneticisidir.
var Store = sessions.NewCookieStore([]byte("cok-gizli-bir-anahtar-12345"))

// init fonksiyonu, bu paket yüklendiğinde otomatik olarak bir kere çalışır.
func init() {
	// GOB'a, oturumda saklayacağımız özel tipimizi tanıtıyoruz.
	// Bu olmadan, oturum yöneticisi CartItem'ı nasıl saklayacağını bilemez.
	gob.Register([]model.CartItem{})
}