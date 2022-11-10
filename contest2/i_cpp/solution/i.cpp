#include <iostream>
#include <vector>
#include <algorithm>
#include <cassert>
#pragma GCC optimize("O3,unroll-loops")
#pragma GCC target("avx2,bmi,bmi2,lzcnt,popcnt")

const int	MODULE = 998244353;
const int   MOD119 = 119;
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

inline void apply_fft(std::vector<int>& poly, int log, bool is_invert_fft) {
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
		if (is_invert_fft) {
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
	if (is_invert_fft) {
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
    for (int i = 0; i < first_poly.size(); ++i) {
        first_poly[i] %= MOD119;
    }
}

inline void find_invert_poly(std::vector<int>& invert, std::vector<int>& poly, int deg) {
    if (poly[0] == 0) {
		return;
	}
	if (poly.size() < deg) {
		poly.resize(deg);
	}
	invert.resize(poly.size());
	int	curDeg = 1;
	
	invert[0] = invert_number(poly[0], MOD119);
	while (curDeg < deg) {
        std::vector<int> firstPart(poly.begin(), poly.begin() + curDeg);
		std::vector<int> secondPart(poly.begin() + curDeg, poly.end());
		multiply_polynoms(firstPart, invert);
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
		multiply_polynoms(secondPart, invert);
		secondPart.resize(deg);
		firstPart.resize(deg);
		for (int i = 0; i < std::max(firstPart.size(), secondPart.size()); ++i) {
			if (i == firstPart.size()) {
                firstPart.push_back(0);
			}
			if (i < firstPart.size() && i < secondPart.size()) {
				firstPart[i] = (-(firstPart[i] % MOD119) - (secondPart[i] % MOD119) + 2*MOD119) % MOD119;
			} else if (i < secondPart.size()) {
				firstPart[i] = (-(secondPart[i] % MOD119) + MOD119) % MOD119;
			} else {
				firstPart[i] = (-(firstPart[i] % MOD119) + MOD119) % MOD119;
			}
		}
		multiply_polynoms(firstPart, invert);
		while (firstPart.size() > deg) {
			firstPart.pop_back();
		}
		for (int i = 0; i < firstPart.size(); ++i) {
			if (i + curDeg == invert.size()) {
                invert.push_back(0);
			}
			invert[i + curDeg] = (invert[i + curDeg] + firstPart[i]) % MOD119;
		}
		curDeg <<= 1;
	}
}

inline void divide_polynoms(std::vector<int>& first_poly, std::vector<int>& second_poly, std::vector<int>& Q, std::vector<int>& R) {
	std::vector<int> second_copy = second_poly;
    int deg = first_poly.size() - second_poly.size() + 1;
	find_invert_poly(Q, second_copy, deg);
	multiply_polynoms(Q, first_poly);
	Q.resize(deg);
    std::reverse(Q.begin(), Q.end());
    std::reverse(first_poly.begin(), first_poly.end());
    std::reverse(second_poly.begin(), second_poly.end());
	multiply_polynoms(second_poly, Q);
	delete_zeros(second_poly);
	for (int i = 0; i < second_poly.size(); ++i) {
		int last = (first_poly[i] + MOD119 - second_poly[i] % MOD119) % MOD119;
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

inline void get_bit_mask(int n, std::vector<int>& mask) {
    while (n) {
        mask.push_back(n % 2);
        n >>= 1;
    }
    std::reverse(mask.begin(), mask.end());
}

inline void fill_char_polynom(input_data& data) {
    char_polynom.resize(data.q + 1);
    char_polynom[0] = (MOD119 * (std::abs(data.b) / MOD119 + 1) - data.b) % MOD119;
    char_polynom[data.q - data.p] = (MOD119 * (std::abs(data.a) / MOD119 + 1) - data.a) % MOD119;
	char_polynom[data.q] = 1;
	assert(char_polynom[data.q - data.p] >= 0);
	assert(char_polynom[0] >= 0);
}

inline void multiply2(std::vector<int>& old_representation, std::vector<int>& new_representation, input_data& data) {
	std::vector<int> e_polynom = old_representation;
    multiply_polynoms(e_polynom, old_representation);
    std::vector<int> quotient;
    std::vector<int> rest;
	std::vector<int> char_polynom_copy = char_polynom;
	delete_zeros(e_polynom);
	delete_zeros(char_polynom_copy);
	if (e_polynom.size() < char_polynom_copy.size()) {
		for (size_t i = 0; i < new_representation.size(); ++i) {
			if (i < e_polynom.size()) {
				new_representation[i] = e_polynom[i] % MOD119;
			} else {
				new_representation[i] = 0;
			}
		}
		return;
	}
	std::reverse(e_polynom.begin(), e_polynom.end());
	std::reverse(char_polynom_copy.begin(), char_polynom_copy.end());
    divide_polynoms(e_polynom, char_polynom_copy, quotient, rest);
    for (size_t i = 0; i < new_representation.size(); ++i) {
		if (i < rest.size()) {
        	new_representation[i] = rest[i] % MOD119;
		} else {
			new_representation[i] = 0;
		}
    }
}

inline void plus1(std::vector<int>& old_represenation, std::vector<int>& new_representation, input_data& data) {
    new_representation[0] = (data.b * 1ll * old_represenation[old_represenation.size() - 1]) % MOD119;
    for (size_t i = 1; i < old_represenation.size(); ++i) {
        new_representation[i] = old_represenation[i - 1];
    }
	new_representation[data.q - data.p] = (new_representation[data.q - data.p] + (data.a * 1ll * old_represenation[old_represenation.size() - 1]) % MOD119) % MOD119;
}

inline void find_base_elements(std::vector<int>& elements, input_data& data) {
    for (int i = 0; i < data.q; ++i) {
        if (i == 0) {
            elements.push_back(1);
            continue;
        }
        int last_element = 0;
        if (i - data.p <= 0) {
            last_element = (last_element + data.a) % MOD119;
        } else {
            last_element = (last_element + (data.a * 1ll * elements[i - data.p]) % MOD119) % MOD119;
        }

        if (i - data.q <= 0) {
            last_element = (last_element + data.b) % MOD119;
        } else {
            last_element = (last_element + (data.b * 1ll * elements[i - data.q]) % MOD119) % MOD119;
        }
        elements.push_back(last_element % MOD119);
    }
}

inline int find_rest(input_data& data) {
	data.a %= MOD119;
	data.b %= MOD119;
	if (data.n < data.q) {
		std::vector<int> base_elements;
		find_base_elements(base_elements, data);
		return base_elements[data.n] % MOD119;
	}
    std::vector<int> mask;
    get_bit_mask(data.n, mask);
	fill_char_polynom(data);
    std::vector<int> cur_representation;
    std::vector<int> new_representation;
    new_representation.resize(data.q);
    cur_representation.resize(data.q);
    cur_representation[0] = data.b % MOD119;
	cur_representation[data.q - data.p] = data.a % MOD119;
	int cur_number = 1;
	bool flag = false;
    for (size_t i = 1; i < mask.size(); ++i) {
		if (!flag)
			cur_number *= 2;
		if (!flag && mask[i] == 1)
			++cur_number;
		if (!flag && cur_number >= data.q) {
			flag = true;
		 	for (int j = 0; j < cur_number - data.q; ++j) {
				plus1(cur_representation, new_representation, data);
				cur_representation = new_representation;
			}
			continue;
		}
		if (!flag) {
			continue;
		}
		assert(cur_representation.size() == new_representation.size());
        multiply2(cur_representation, new_representation, data);
		cur_representation = new_representation;
        if (mask[i] == 1) {
            plus1(cur_representation, new_representation, data);
			cur_representation = new_representation;
        }
    }
    std::vector<int> base_elements;
    find_base_elements(base_elements, data);
    int ans = 0;
    for (size_t i = 0; i < cur_representation.size(); ++i) {
        ans = (ans + (cur_representation[i] * 1ll * base_elements[i]) % MOD119) % MOD119;
    }
    return ans % MOD119;
}

// int main() {
//     std::ios_base::sync_with_stdio(false);
// 	std::cin.sync_with_stdio(false);
// 	std::cout.sync_with_stdio(false);
//     std::cin.tie(NULL);
//     std::cout.tie(NULL);
//     input_data data;
//     std::cin >> data.n >> data.a >> data.b >> data.p >> data.q;
//     int ans = find_rest(data);
//     std::cout << ans << '\n';
// }