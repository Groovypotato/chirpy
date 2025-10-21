package main

func validateChirpHandler(msg string) (bool, int, string) {
	if len(msg) > 140 {
		return false, 400, "Chirp is too long"
	} else {
		chirp := badWordFilterHandler(msg)
		return true, 255, chirp
	}
}
