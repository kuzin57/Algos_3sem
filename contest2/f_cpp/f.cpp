#include <iostream>
#include <vector>
#include <algorithm>
#pragma GCC optimize("O3,unroll-loops")
#pragma GCC target("avx2,bmi,bmi2,lzcnt,popcnt")

const int	MODULE = 7340033;
const int	ROOT   = 5;
const int	INVERT = 4404020;
const int	POWER  = 1 << 20;
inline std::vector<int> rev;

inline int getReversed(int num, int log) {
	int left = 0;
	int right = log - 1;
	while (left < right) {
		int bitLeft = (num & (1 << left)) >> left;
		int bitRight = (num & (1 << right)) >> right;
		num ^= (bitLeft << left);
		num ^= (bitRight << left);
		num ^= (bitRight << right);
		num ^= (bitLeft << right);
		++left;
		--right;
	}
	return num;
}

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

inline int invert_number(int num) {
	int x, y;
	gcd(num, MODULE, x, y);
	while (x < 0) {
		x = (x + MODULE) % MODULE;
	}
	return x;
}

inline void fill_rev(int n, int log) {
	for (int i = 0; i < n; ++i) {
		rev[i] = getReversed(i, log);
	}
}

inline void delete_zeros(std::vector<int>& poly) {
	while (poly.size() > 1 && poly[poly.size() - 1] == 0) {
		poly.pop_back();
	}
}

inline void fft(std::vector<int>& poly, int log, bool isInvertFFT) {
	if (poly.size() == 1) {
		return;
	}

	for (int i = 0; i < poly.size(); ++i) {
		int reversed = rev[i];
		if (i < reversed) {
			std::swap(poly[i], poly[reversed]);
		}
	}

	int curOffset = 1;
	for (int i = 0; i < log; ++i) {
		int root = 0;
		if (isInvertFFT) {
			root = INVERT;
		} else {
			root = ROOT;
		}
		for (int j = curOffset << 1; j < POWER; j <<= 1) {
			root = ((root * 1ll) * root) % MODULE;
		}
		for (int j = 0; j < poly.size(); j += curOffset << 1) {
			int curRoot = 1;
			for (int k = 0; k < curOffset; ++k) {
				int tmp1 = poly[k + j];
				int tmp2 = poly[k + j + curOffset];
				poly[k + j] = tmp1 + ((curRoot*1ll)*tmp2)%MODULE;
				if (poly[k + j] > MODULE) {
					poly[k + j] -= MODULE;
				}
				poly[k + j + curOffset] = tmp1 - ((curRoot*1ll)*tmp2)%MODULE;
				if (poly[k + j + curOffset] < 0) {
					poly[k + j + curOffset] += MODULE;
				}
				curRoot = ((curRoot * 1ll) * root) % MODULE;
			}
		}
		curOffset <<= 1;
	}
	if (isInvertFFT) {
		int inverted = invert_number(poly.size());
		for (int i = 0; i < poly.size(); ++i) {
			poly[i] = ((poly[i] * 1ll) * inverted) % MODULE;
		}
	}
}

inline void multiply(std::vector<int>& first_poly, std::vector<int>& second_poly_arg) {	
	if (first_poly.size() == 0 || second_poly_arg.size() == 0) {
		first_poly.resize(1);
		first_poly[0] = 0;
		return;
	}
	int first_size = first_poly.size();
	while (first_size > 1 && first_poly[first_size-1] == 0) {
        --first_size;
	}
	int second_size = second_poly_arg.size();
	while (second_size > 1 && second_poly_arg[second_size-1] == 0) {
		--second_size;
	}

	if (first_poly[first_size-1] == 0 || second_poly_arg[second_size-1] == 0) {
		first_poly.resize(1);
		first_poly[0] = 0;
		return;
	}

	std::vector<int> second_poly = second_poly_arg;
	int min_deg_two = 1;
	int log = 0;
	while (min_deg_two < std::max(first_size, second_size)) {
		++log;
		min_deg_two <<= 1;
	}
	min_deg_two <<= 1;
	++log;
	rev.resize(min_deg_two);
	fill_rev(min_deg_two, log);

	first_poly.resize(min_deg_two);
	second_poly.resize(min_deg_two);

	fft(first_poly, log, false);
	fft(second_poly, log, false);

	for (int i = 0; i < first_poly.size(); ++i) {
		first_poly[i] = ((first_poly[i] * 1ll) * second_poly[i]) % MODULE;
	}

	fft(first_poly, log, true);
}

inline void findInvertPoly(std::vector<int>& invert, std::vector<int>& poly, int deg) {
    if (poly[0] == 0) {
		return;
	}
	if (poly.size() < deg) {
		poly.resize(deg);
	}
	invert.resize(poly.size());
	int	curDeg = 1;
	
	invert[0] = invert_number(poly[0]);
	while (curDeg < deg) {
        std::vector<int> firstPart(poly.begin(), poly.begin() + curDeg);
		std::vector<int> secondPart(poly.begin() + curDeg, poly.end());
		multiply(firstPart, invert);
		while (firstPart.size() > deg) {
			firstPart.pop_back();
		}
		if (curDeg >= firstPart.size()) {
			firstPart.resize(1);
			firstPart[0] = 0;
		} else {
			std::reverse(firstPart.begin(), firstPart.end());
			for (int i = 0; i < curDeg; ++i) {
				firstPart.pop_back();
			}
			std::reverse(firstPart.begin(), firstPart.end());
		}
		multiply(secondPart, invert);
		while (secondPart.size() > deg) {
			secondPart.pop_back();
		}
		for (int i = 0; i < std::max(firstPart.size(), secondPart.size()); ++i) {
			if (i == firstPart.size()) {
                firstPart.push_back(0);
			}
			if (i < firstPart.size() && i < secondPart.size()) {
				firstPart[i] = (-firstPart[i] - secondPart[i] + 2*MODULE) % MODULE;
			} else if (i < secondPart.size()) {
				firstPart[i] = (-secondPart[i] + MODULE) % MODULE;
			} else {
				firstPart[i] = (-firstPart[i] + MODULE) % MODULE;
			}
		}
		multiply(firstPart, invert);
		while (firstPart.size() > deg) {
			firstPart.pop_back();
		}
		for (int i = 0; i < firstPart.size(); ++i) {
			if (i + curDeg == invert.size()) {
                invert.push_back(0);
			}
			invert[i + curDeg] = (invert[i + curDeg] + firstPart[i]) % MODULE;
		}
		curDeg <<= 1;
	}
}

int main() {
    std::ios_base::sync_with_stdio(false);
	std::cin.sync_with_stdio(false);
	std::cout.sync_with_stdio(false);
    std::cin.tie(NULL);
    std::cout.tie(NULL);
    int m;
    std::cin >> m;
    int n;
    std::cin >> n;
    std::vector<int> poly(n + 1);
	for (int i = 0; i <= n; ++i) {
        std::cin >> poly[i];
	}
	std::vector<int> result;
	findInvertPoly(result, poly, m);
	if (result.size() == 0) {
		std::cout << "The ears of a dead donkey" << '\n';
	} else {
		result.resize(m);
		for (int i = 0; i < result.size(); ++i) {
			std::cout << (result[i] % MODULE) << " ";
		}
		std::cout << '\n';
	}
}