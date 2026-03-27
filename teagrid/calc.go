package teagrid

func gcd(x, y int) (result int) {
	if x == 0 {
		result = y
		goto end
	}
	if y == 0 {
		result = x
		goto end
	}
	result = gcd(y%x, x)
end:
	return result
}
