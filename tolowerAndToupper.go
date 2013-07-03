逻辑或
func tolower(c byte) byte {
  return c | 0x20
}

逻辑异或（不一样为真）
func toupper(c byte) byte {
	return c ^ 0x20
}
