#include <iostream>

inline int gcd (int a, int b, int& x, int& y) {
	if (a == 0) {
		x = 0; 
		y = 1;
		return b;
	}
	int x_new, y_new;
	int d = gcd(b%a, a, x_new, y_new);
	x = (y_new - ((b / a) * 1ll) * x_new);
	y = x_new;
	return d;
}

inline int invert_number(int num, int mod) {
	int x, y;
	gcd(num, mod, x, y);
	while (x < 0) {
		x = (x + mod) % mod;
	}
	return x;
}

int main() {
    std::cout << invert_number(31, 998244353) << '\n';
}