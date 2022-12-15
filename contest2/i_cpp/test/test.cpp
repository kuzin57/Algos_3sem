#include <gtest/gtest.h>
#include <gmock/gmock.h>
#include <chrono>
#include <random>
#include <fstream>
#include "../solution/i.cpp"

const int NAB_LIMIT = 2000000000;
const int PQ_LIMIT = 10000;

void gen_test(input_data& data) {
    std::random_device dev;
    std::mt19937 rng(dev());
    std::uniform_int_distribution<std::mt19937::result_type> dist(1,10000); 
    std::uniform_int_distribution<std::mt19937::result_type> dist2(0,2000000000);
    std::uniform_int_distribution<std::mt19937::result_type> dist3(1,20000000);
    data.n = dist3(rng);

    data.a = dist2(rng);
    data.b = dist2(rng);

    data.p = dist(rng);

    data.q = dist(rng);
    if (data.q <= data.p) {
        data.q = data.p + 1;
    }
}

int stupid_algo(input_data& data) {
    std::vector<int> sequence;
    sequence.push_back(1);
    for (int i = 1; i <= data.n; ++i) {
        int last_element = 0;
        if (i - data.p <= 0) {
            last_element = (last_element + data.a) % MOD119;
        } else {
            last_element = (last_element + sequence[i - data.p] * 1ll * data.a % MOD119) % MOD119;
        }

        if (i - data.q <= 0) {
            last_element = (last_element + data.b) % MOD119;
        } else {
            last_element = (last_element + sequence[i - data.q] * 1ll * data.b % MOD119) % MOD119;
        }
        sequence.push_back(last_element % MOD119);
    }
    return sequence[data.n];
}

void print_test(input_data& data) {
    std::cout << data.n << " " << data.a << " " << data.b << " " << data.p << " " << data.q << std::endl;
}

TEST(TestSequence, TestCorrect) {
    for (int i = 0; i < 100; ++i) {
        rev.clear();
        char_polynom.clear();
        input_data data;
        gen_test(data);
        input_data copy = data;
        int answer_to_check = find_rest(data);
        int right_answer = stupid_algo(copy);
        if (answer_to_check != right_answer) {
            print_test(data);
        }
        ASSERT_EQ(answer_to_check, right_answer);
    }
}

int main(int argc, char **argv)
{
    ::testing::InitGoogleTest(&argc, argv);
    ::testing::InitGoogleMock(&argc, argv);
    
    return RUN_ALL_TESTS();
}