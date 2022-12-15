#include <iostream>
#include <vector>
#include <algorithm>
#include <cassert>
#pragma GCC optimize("O3,unroll-loops")
#pragma GCC target("avx2,bmi,bmi2,lzcnt,popcnt")

const int	MODULE = 998244353;
const int	ROOT   = 31;
const int	INVERT = 128805723;
const int	POWER  = 1 << 23;
inline std::vector<int> rev;
inline std::vector<int> char_polynom;

struct input_data {
    int n;
    int p;
    int q;
    int a;
    int b;
};

inline int get_reversed(int num, int log) {
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
		rev[i] = get_reversed(i, log);
	}
}

inline void delete_zeros(std::vector<int>& poly) {
	while (poly.size() > 1 && poly[poly.size() - 1] == 0) {
		poly.pop_back();
	}
}

inline void apply_fft(std::vector<int>& poly, int log, bool isInvertapply_fft) {
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
		if (isInvertapply_fft) {
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
	if (isInvertapply_fft) {
		int inverted = invert_number(poly.size(), MODULE);
		for (int i = 0; i < poly.size(); ++i) {
			poly[i] = ((poly[i] * 1ll) * inverted) % MODULE;
		}
	}
}

inline void multiply_polynoms(std::vector<int>& first_poly, std::vector<int>& second_poly_arg) {	
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

	apply_fft(first_poly, log, false);
	apply_fft(second_poly, log, false);

	for (int i = 0; i < first_poly.size(); ++i) {
		first_poly[i] = ((first_poly[i] * 1ll) * second_poly[i]) % MODULE;
	}

	apply_fft(first_poly, log, true);
    delete_zeros(first_poly);
}

inline void get_bit_mask(int n, std::vector<int>& mask) {
    while (n) {
        mask.push_back(n % 2);
        n >>= 1;
    }
    std::reverse(mask.begin(), mask.end());
}

inline void sum_polynoms(std::vector<int>& first_poly, std::vector<int>& second_poly) {
    for (size_t i = 0; i < std::max(first_poly.size(), second_poly.size()); ++i) {
        if (i < first_poly.size() && i < second_poly.size()) {
            first_poly[i] = (first_poly[i] + second_poly[i]) % MODULE;
            continue;
        }
        if (i < second_poly.size()) {
            first_poly.push_back(second_poly[i] % MODULE);
            continue;
        }
    }
}

inline int find_ways(int n, int k) {
    std::vector<int> without_first(k + 1);
    std::vector<int> without_last(k + 1);
    std::vector<int> without_first_last(k + 1);
    std::vector<int> with_first_last(k + 1);
    std::vector<int> sum_polynom(k + 1);
    without_first[0] = 1;
    without_first[1] = 1;
    without_first[2] = 0;
    without_last[0] = 1;
    without_last[1] = 1;
    without_last[2] = 0;
    without_first_last[0] = 1;
    without_first_last[1] = 0;
    without_first_last[2] = 0;
    with_first_last[0] = 1;
    with_first_last[1] = 1;
    with_first_last[2] = 2;
    sum_polynom[0] = 1;
    sum_polynom[1] = 2;
    sum_polynom[2] = 2;
    std::vector<int> mask;
    get_bit_mask(n, mask);
    for (size_t i = 1; i < mask.size(); ++i) {
        std::vector<int> tmp1 = without_first;
        multiply_polynoms(tmp1, sum_polynom);
        std::vector<int> tmp2 = without_first_last;
        multiply_polynoms(tmp2, sum_polynom);
        std::vector<int> tmp3 = without_first_last;
        multiply_polynoms(tmp3, without_first);

    }
}

int main() {
    int n, k;
    std::cin >> n >> k;
}