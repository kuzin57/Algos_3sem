#include <iostream>
#include <vector>
#include <cassert>

const int N = 10000000;

void get_primes(std::vector<int>& primes, std::vector<int>& mind) {
    for (int i = 2; i <= N; ++i) {
        if (mind[i] == i)
            primes.push_back(i);
        for (size_t j = 0; j < primes.size(); ++j) {
            if (primes[j] * i >= mind.size() || primes[j] > mind[i])
                break;
            mind[i*primes[j]] = primes[j];
        }
    }
}

int get_min_not_used_prime(int index, const std::vector<int>& primes, std::vector<bool>& primesUsed) {
    for (int i = index; i < primes.size(); ++i) {
        if (!primesUsed[i]) {
            primesUsed[i] = true;
            return i;
        }
    }
    return -1;
}

int get_first_number(int num, const std::vector<int>& primes, std::vector<bool>& primesUsed) {
    bool exit;
    while (true) {
        exit = true;
        for (size_t i = 0; i < primes.size(); ++i) {
            if (primes[i] > num)
                break;
            if (num % primes[i] == 0 && primesUsed[i]) {
                exit = false;
                break;
            }
        }
        if (exit) 
            break;
        ++num;
    }
    return num;
}

int main() {
    int n;
    int index = 0;
    bool take_only_primes = false;
    std::cin >> n;
    std::vector<int> arr(n);
    std::vector<int> mind(N + 1);
    for (int i = 2; i <= N; ++i) {
        mind[i] = i;
    }
    std::vector<int> primes;
    get_primes(primes, mind);
    std::vector<bool> primes_used(primes.size());
    for (int i = 0; i < n; ++i) {
        std::cin >> arr[i];
    }
    for (int i = 0; i < n; ++i) {
        if (take_only_primes) {
            index = get_min_not_used_prime(index, primes, primes_used);
            arr[i] = primes[index];
            continue;
        }
        for (size_t j = 0; j < primes.size(); ++j) {
            if (arr[i] < primes[j]) {
                break;
            }
            if (arr[i] % primes[j] == 0 && primes_used[j]) {
                take_only_primes = true;
                arr[i] = get_first_number(arr[i], primes, primes_used);
                break;
            }
        }
        for (size_t j = 0; j < primes.size(); ++j) {
            if (arr[i] < primes[j]) 
                break;
            if (arr[i] % primes[j] == 0)
                primes_used[j] = true;
        }
    }
    for (size_t i = 0; i < arr.size(); ++i) {
        std::cout << arr[i] << " ";
    }
    std::cout << '\n';
}