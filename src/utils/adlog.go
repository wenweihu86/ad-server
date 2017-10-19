package utils
import(
  "os"
)

func CheckFileIsExist(filename string) (bool) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false;
	}
	return true;
}
