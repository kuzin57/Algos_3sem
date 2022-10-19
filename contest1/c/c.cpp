#include <iostream>
#include <vector>
#include <unordered_map>
#include <string>
#include <memory>
#include <map>
#include <algorithm>

struct State {
    State();
    State(int, int, char, int, int);
    std::map<int, int> to;
    std::map<int, int> go;
    std::vector<int> word_number;
    int depth;
    int parent;
    int link = -1;
    int compressed_link = -1;
    char letter;
};

State::State() {}

State::State(int depth, int parent, char letter, int link, int compressed_link)
    : depth(depth), parent(parent), letter(letter), link(link), compressed_link(compressed_link) {}

struct Bohr {
    Bohr() = default;
    std::vector<std::shared_ptr<State>> states;
    std::map<char, int> alphabet;
    int start_state;

    void set_root();
    void initialize_alphabet(const std::string& line, const std::vector<std::string>& dict);
    void add(const std::string& word, int index);
    int get_link(int s);
    int go(int s, char ch);
    int get_compressed_link(int s);
    void find_occurences(std::vector<std::vector<int>>& occurences, const std::string& line, int size);
};

void Bohr::initialize_alphabet(const std::string& line, const std::vector<std::string>& dict) {
    int cnt = 0;
    for (int i = 0; i < line.size(); ++i) {
        auto it = alphabet.find(line[i]);
        if (it == alphabet.end()) {
            alphabet[line[i]] = cnt++;
        }
    }

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
    new_state->link = start_state;
    new_state->depth = 0;
    states.push_back(new_state);
}

void Bohr::add(const std::string& word, int index) {
    auto cur_state = states[start_state];
    int cur_index = 0;
    for (int i = 0; i < word.size(); ++i) {
        auto next = cur_state->to.find(alphabet[word[i]]);
        if (next == cur_state->to.end()) {
            auto new_state = std::make_shared<State>(cur_state->depth + 1, cur_index, word[i], -1, -1);
            states.push_back(new_state);
            cur_state->to[alphabet[word[i]]] = states.size() - 1;
            next->second = states.size() - 1;
        } 
        cur_state = states[next->second];
        cur_index = next->second;
    }
    cur_state->word_number.push_back(index);
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

int Bohr::get_compressed_link(int s) {
    auto st = states[s];
    if (st->link == -1) {
        st->link = get_link(s);
    }
    auto st_link = states[st->link];
    if (st->compressed_link == -1) {
        if (st->parent == start_state || s == start_state) {
            st->compressed_link = start_state;
        } else if (st_link->word_number.size() > 0) {
            st->compressed_link = st->link;
        } else {
            st->compressed_link = get_compressed_link(st->link);
        }
    }
    return st->compressed_link;
}

void Bohr::find_occurences(std::vector<std::vector<int>>& occurences, const std::string& line, int size) {
    auto cur = start_state;
    for (size_t i = 0; i < line.size(); ++i) {
        auto ss = states[cur];
        auto next = go(cur, line[i]);
        auto node = next;
        while (node != start_state) {
            auto n = states[node];
            if (n->word_number.size() > 0) {
                for (size_t j = 0; j < n->word_number.size(); ++j) {
                    occurences[n->word_number[j]].push_back(i + 2 - n->depth);
                }
            }
            node = get_compressed_link(node);
        }
        cur = next;
    }
}

void run(const std::string& line, const std::vector<std::string>& dictionary, std::vector<std::vector<int>>& occurences) {
    auto bohr = Bohr();
    bohr.initialize_alphabet(line, dictionary);
    bohr.set_root();
    for(int i = 0; i < dictionary.size(); ++i) {
        bohr.add(dictionary[i], i);
    }
    bohr.find_occurences(occurences, line, dictionary.size());
}

int main() {
    std::string line;
    std::string word;
    int N;
    std::cin >> line;

    std::cin >> N;
    auto dictionary = std::vector<std::string>(N);
    auto occurences = std::vector<std::vector<int>>(N);
    for(int i = 0; i < N; ++i) {
        std::cin >> dictionary[i];
    }

    run(line, dictionary, occurences);
    for(size_t i = 0; i < occurences.size(); ++i) {
        std::cout << occurences[i].size() << " ";
        for(size_t j = 0; j < occurences[i].size(); ++j) {
            std::cout << occurences[i][j] << " ";
        }
        std::cout << std::endl;
    }
}