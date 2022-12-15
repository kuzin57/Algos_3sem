#include <gtest/gtest.h>
#include <gmock/gmock.h>
#include <chrono>
#include <random>
#include <fstream>
#include "../solution/g.cpp"

void gen_test(std::vector<int>& first, std::vector<int>& second, int length) {
    std::random_device dev;
    std::mt19937 rng(dev());
    std::uniform_int_distribution<std::mt19937::result_type> coefficients_generator(0,6);
    std::uniform_int_distribution<std::mt19937::result_type> length_generator(1,length);

    int n = length_generator(rng) + 1;
    for (int i = 0; i < n + 1; ++i) {
        int coeff = coefficients_generator(rng);
        first.push_back(coeff);
    }

    std::uniform_int_distribution<std::mt19937::result_type> length_generator_second(1,n / 2);
    int m = rand() % (n / 2) + 1;
    for (int i = 0; i < m + 1; ++i) {
        int coeff = rand() % 7;
        second.push_back(coeff);
    }
}

void print_test(std::ofstream& out, std::vector<int>& F, std::vector<int>& G) {
    std::reverse(F.begin(), F.end());
    std::reverse(G.begin(), G.end());
    out << F.size() - 1 << " ";
    for (int i = 0; i < F.size(); ++i) {
        out << F[i] << " ";
    }
    out << std::endl;
    out << G.size() - 1 << " ";
    for (int i = 0; i < G.size(); ++i) {
        out << G[i] << " ";
    }
    out << std::endl;
}

bool check(std::vector<int>& F, std::vector<int>& G, std::vector<int>& Q, std::vector<int>& R) {
    std::ofstream failed_test_output("test.txt", std::fstream::out);
    for (size_t i = 0; i < Q.size(); ++i) {
        if (!(Q[i] >= 0 && Q[i] <= 6)) {
            std::cout << "problem with module in Q" << std::endl;
            print_test(failed_test_output, F, G);
            return false;
        }
    }
    for (size_t i = 0; i < R.size(); ++i) {
        if (!(R[i] >= 0 && R[i] <= 6)) {
            std::cout << "problem with module in R" << std::endl;
            print_test(failed_test_output, F, G);
            return false;
        }
    }
    if ((Q[Q.size() - 1] == 0 && Q.size() > 1) || (R[R.size() - 1] == 0 && R.size() > 1)) {
        std::cout << "wrong format" << std::endl;
        print_test(failed_test_output, F, G);
        return false;
    }
    int GSize = G.size();
    multiply(G, Q);
    if (G.size() != F.size()) {
        std::cout << "size problem: " << F.size() << " " << GSize << " " << Q.size() << " " << R.size() << std::endl;
        print_test(failed_test_output, F, G);
        return false;
    }
    for (size_t i = 0; i < F.size(); ++i) {
        if (i == G.size()) {
            G.push_back(0);
        }
        if (i == R.size()) {
            R.push_back(0);
        }
        if (!((F[i] - G[i] - R[i]) % 7 == 0)) {
            std::cout << "ooops" << std::endl;
            print_test(failed_test_output, F, G);
            return false;
        }
    }
    return true;
}

TEST(TestDivPoly, testCorrect) {
    for (int i = 0; i < 1000; ++i) {
        std::vector<int> first;
        std::vector<int> second;
        gen_test(first, second, 50000);
        if (first.size() < second.size()) {
            continue;
        }
        if (first[0] == 0) 
            continue;
        if (second[0] == 0)
            continue;
        std::vector<int> Q;
        std::vector<int> R;
        std::vector<int> F = first;
        std::vector<int> G = second;
        std::vector<int> first_copy = first;
        std::vector<int> second_copy = second;
        divPoly(F, G, Q, R);
        std::reverse(first.begin(), first.end());
        std::reverse(second.begin(), second.end());
        ASSERT_TRUE(check(first, second, Q, R));
    }
}

TEST(TestDivPoly, performanceTests) {
    for (int i = 0; i < 10; ++i) {
        std::vector<int> first;
        std::vector<int> second;
        gen_test(first, second, 50000);
        if (first.size() < second.size()) {
            continue;
        }
        if (first[0] == 0) 
            continue;
        if (second[0] == 0)
            continue;
        std::vector<int> Q;
        std::vector<int> R;
        std::vector<int> F = first;
        std::vector<int> G = second;
        auto start = std::chrono::high_resolution_clock::now();
        divPoly(F, G, Q, R);
        auto stop = std::chrono::high_resolution_clock::now();
        auto duration = std::chrono::duration_cast<std::chrono::milliseconds>(stop - start);
        int64_t t = duration.count();
        std::cout << "duration: " << duration.count() << " millieseconds" << '\n';
        ASSERT_TRUE(t <= 1000);
    }
}

int main(int argc, char **argv)
{
  ::testing::InitGoogleTest(&argc, argv);
  ::testing::InitGoogleMock(&argc, argv);
  
  return RUN_ALL_TESTS();
}