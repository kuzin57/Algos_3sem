#include <iostream>
#include <vector>
#include <unordered_map>
#include <string>
#include <memory>
#include <map>
#include <algorithm>
#include <cassert>

const int ALPHABET_SIZE = 26;
const int MODULE = 10000;

struct State {
    State();
    State(int, int, char, int);
    std::map<int, int> to;
    std::map<int, int> go;
    int depth = -1;
    int parent = -1;
    int link = -1;
    char letter;
    bool is_term = false;
};

State::State() {}

State::State(int depth, int parent, char letter, int link)
    : depth(depth), parent(parent), letter(letter), link(link) {}

struct Bohr {
    Bohr() = default;
    std::vector<std::shared_ptr<State>> states;
    std::map<char, int> alphabet;
    std::vector<int> dp;
    int start_state;

    void set_root();
    void initialize_alphabet(const std::vector<std::string>& dict);
    void add(const std::string& word, int index);
    int get_link(int s);
    int go(int s, char ch);
    void count(int index, int len, std::vector<int>& new_dp);
};

void Bohr::initialize_alphabet(const std::vector<std::string>& dict) {
    int cnt = 0;
    for(size_t i = 0; i < dict.size(); ++i) {
        for (size_t j = 0; j < dict[i].size(); ++j) {
            auto it = alphabet.find(dict[i][j]);
            if (it == alphabet.end()) {
                alphabet[dict[i][j]] = cnt++;
            }
        }
    }
}

void Bohr::set_root() {
    start_state = 0;
    auto new_state = std::make_shared<State>();
    new_state->depth = 0;
    states.push_back(new_state);
}

void Bohr::add(const std::string& word, int index) {
    auto cur_state = states[start_state];
    int cur_index = 0;
    for (int i = 0; i < word.size(); ++i) {
        auto next = cur_state->to.find(alphabet[word[i]]);
        if (next == cur_state->to.end()) {
            auto new_state = std::make_shared<State>(cur_state->depth + 1, cur_index, word[i], -1);
            states.push_back(new_state);
            cur_state->to[alphabet[word[i]]] = states.size() - 1;
            next->second = states.size() - 1;
        } 
        cur_state = states[next->second];
        cur_index = next->second;
    }
    cur_state->is_term = true;
}

int Bohr::get_link(int s) {
    auto st = states[s];
    if (st->link == -1) {
        if (st->parent == start_state || s == start_state) {
            st->link = start_state;
        } else {
            st->link = go(get_link(st->parent), st->letter);
        }
    }
    return st->link;
}
 
int Bohr::go(int s, char ch) {
    if (alphabet.find(ch) == alphabet.end()) {
        return start_state;
    }
    int c = alphabet[ch];
    auto st = states[s];
    auto next = st->go.find(c);
    if (next == st->go.end()) {
        auto n = st->to.find(c);
        if (n != st->to.end()) {
            st->go[c] = st->to[c];
        } else {
            if (s == start_state) {
                st->go[c] = start_state;
            } else {
                st->go[c] = go(get_link(s), ch);
            }
        }
    }
    return st->go[c];
}

void Bohr::count(int index, int len, std::vector<int>& new_dp) {
    if (len < states[index]->depth || states[index]->is_term || states[get_link(index)]->is_term) {
        return;
    }
 
    for (int i = 0; i < ALPHABET_SIZE; ++i) {
        int golink = go(index, 'a'+i);
        if (states[get_link(golink)]->is_term) {
            continue;
        }
        new_dp[golink] += dp[index];
        new_dp[golink] %= MODULE;
    }

    for (const auto& next : states[index]->to) {
        count(next.second, len, new_dp);
    }
}

int run(int length, const std::vector<std::string>& dictionary) {
    auto bohr = Bohr();
    bohr.initialize_alphabet(dictionary);
    bohr.set_root();
    for(int i = 0; i < dictionary.size(); ++i) {
        bohr.add(dictionary[i], i);
    }
    bohr.dp = std::vector<int>(bohr.states.size());
    bohr.dp[bohr.start_state] = 1;

    for (int i = 0; i < length; ++i) {
        std::vector<int> new_dp(bohr.states.size());
        bohr.count(bohr.start_state, i + 1, new_dp);
        for (size_t j = 0; j < bohr.dp.size(); ++j) {
            bohr.dp[j] = new_dp[j];
        }
    }

    int ans = 0;
    for (int i = 0; i < bohr.states.size(); ++i) {
        if (!bohr.states[i]->is_term) {
            ans = (ans + bohr.dp[i]) % MODULE;
        }
    }

    int rest = 1;
    for (int i = 0; i < length; ++i) {
        rest = (rest * ALPHABET_SIZE) % MODULE;
    }
    
    if (rest < ans) {
        rest += MODULE;
    }
    return (rest - ans) % MODULE;
}

int stupid_algo(int n, int k, std::vector<std::string>& dictionary) {
    std::vector<int> counters(n);
    int power26 = 1;
    for (int i = 0; i < n; ++i) {
        power26 *= ALPHABET_SIZE;
    }
    int cur_index = 0;
    std::vector<std::string> words;
    for (int j = 0; j < ALPHABET_SIZE; ++j) {
        for (int k = 0; k < ALPHABET_SIZE; ++k) {
            for (int h = 0; h < ALPHABET_SIZE; ++h) {
                for (int l = 0; l < ALPHABET_SIZE; ++l) {
                    std::string w = "";
                    w.push_back('a' + j);
                    w.push_back('a' + k);
                    w.push_back('a' + h);
                    w.push_back('a' + l);
                    words.push_back(w);
                }
            }
        }
    }
    int counter = 0;
    bool exit = false;
    for (size_t t = 0; t < words.size(); ++t) {
        exit = false;
        for (size_t i = 0; i < dictionary.size(); ++i) {
            for (int j = 0; j <= n - dictionary[i].size(); ++j) {
                if (words[t].substr(j, dictionary[i].size()) == dictionary[i]) {
                    ++counter;
                    counter %= MODULE;
                    exit = true;
                    break;
                }
            }
            if (exit) {
                break;
            }
        }
    }
    return counter;
}

int main() {
    int n, k;
    std::cin >> n >> k;

    std::vector<std::string> dictionary(k);
    for(int i = 0; i < k; ++i) {
        std::cin >> dictionary[i];
    }

    int ans = run(n, dictionary);
    std::cout << ans << std::endl;
}