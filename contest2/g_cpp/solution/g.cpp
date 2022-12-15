#include <iostream>
#include <vector>
#include <algorithm>
#include <cassert>
#pragma GCC optimize("O3,unroll-loops")
#pragma GCC target("avx2,bmi,bmi2,lzcnt,popcnt")

const int	MODULE = 7340033;
const int	ROOT   = 5;
const int	INVERT = 4404020;
const int	POWER  = 1 << 20;
const int   MOD7   = 7;
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
	x = (y_new - (b / a * 1ll * x_new));
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
		int inverted = invert_number(poly.size(), MODULE);
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
    delete_zeros(first_poly);
    for (int i = 0; i < first_poly.size(); ++i) {
        first_poly[i] %= MOD7;
    }
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
	
	invert[0] = invert_number(poly[0], MOD7);
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
		secondPart.resize(deg);
		firstPart.resize(deg);
		for (int i = 0; i < std::max(firstPart.size(), secondPart.size()); ++i) {
			if (i == firstPart.size()) {
                firstPart.push_back(0);
			}
			if (i < firstPart.size() && i < secondPart.size()) {
				firstPart[i] = (-(firstPart[i] % MOD7) - (secondPart[i] % MOD7) + 2*MOD7) % MOD7;
			} else if (i < secondPart.size()) {
				firstPart[i] = (-(secondPart[i] % MOD7) + MOD7) % MOD7;
			} else {
				firstPart[i] = (-(firstPart[i] % MOD7) + MOD7) % MOD7;
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
			invert[i + curDeg] = (invert[i + curDeg] + firstPart[i]) % MOD7;
		}
		curDeg <<= 1;
	}
}

inline void divPoly(std::vector<int>& first_poly, std::vector<int>& second_poly, std::vector<int>& Q, std::vector<int>& R) {
	std::vector<int> second_copy = second_poly;
    int deg = first_poly.size() - second_poly.size() + 1;
	findInvertPoly(Q, second_copy, deg);
	multiply(Q, first_poly);
	Q.resize(deg);
    std::reverse(Q.begin(), Q.end());
    std::reverse(first_poly.begin(), first_poly.end());
    std::reverse(second_poly.begin(), second_poly.end());
	multiply(second_poly, Q);
	delete_zeros(second_poly);
	for (int i = 0; i < second_poly.size(); ++i) {
		int last = (first_poly[i] + MOD7 - second_poly[i] % MOD7) % MOD7;
        R.push_back(last);
	}
    delete_zeros(Q);
    delete_zeros(R);
	if (Q.size() == 0) {
		Q.push_back(0);
	}
	if (R.size() == 0) {
		R.push_back(0);
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
    std::vector<int> first_poly(m + 1);
    for (int i = 0; i < m + 1; ++i) {
        std::cin >> first_poly[i];
    }
    int n;
    std::cin >> n;
    std::vector<int> second_poly(n + 1);
	for (int i = 0; i <= n; ++i) {
        std::cin >> second_poly[i];
	}
    if (m < n) {
        std::cout << "0 0" << '\n';
        std::cout << m << " ";
        for (int i = 0; i < m + 1; ++i) {
            std::cout << first_poly[i] << " ";
        }
        std::cout << '\n';
        return 0;
    }
	std::vector<int> Q;
    std::vector<int> R;
	divPoly(first_poly, second_poly, Q, R);
	std::cout << Q.size() - 1 << " ";
    for (int i = Q.size() - 1; i >= 0; --i) {
        std::cout << Q[i] << " ";
    }
    std::cout << '\n';
    std::cout << R.size() - 1 << " ";
    for (int i = R.size() - 1; i >= 0; --i) {
        std::cout << R[i] << " ";
    }
    std::cout << '\n';
}