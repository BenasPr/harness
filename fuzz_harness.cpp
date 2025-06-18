#include <iostream>
#include <vector>
#include <string>
#include <unistd.h>
#include <arpa/inet.h>
#include <sys/socket.h>
#include <netinet/in.h>
#include <fstream>
#include <filesystem>
#include <ctime>

int main(int argc, char* argv[]) {
    std::vector<char> inputBuf((std::istreambuf_iterator<char>(std::cin)),
                                std::istreambuf_iterator<char>());

    if (inputBuf.empty()) {
        std::cerr << "Empty input\n";
        return 1;
    }

    std::string stateFile = "connection_state.flag";
    bool wasConnectedBefore = false;
    if (std::ifstream in{stateFile}) {
        char c;
        in >> c;
        wasConnectedBefore = (c == '1');
    }

    int sock = socket(AF_INET, SOCK_STREAM, 0);
    if (sock < 0) {
        perror("socket");
        return 1;
    }

    sockaddr_in serv_addr{};
    serv_addr.sin_family = AF_INET;
    serv_addr.sin_port = htons(9000);
    if (inet_pton(AF_INET, "192.168.0.146", &serv_addr.sin_addr) <= 0) {
        std::cerr << "Invalid address\n";
        close(sock);
        return 1;
    }

    if (connect(sock, (sockaddr*)&serv_addr, sizeof(serv_addr)) < 0) {
        std::ofstream(stateFile) << "0";
        system("./start_pipeline.sh");
        if (wasConnectedBefore) {
            std::filesystem::create_directory("crashes");
            std::string crashFile = "crashes/crash_" + std::to_string(time(nullptr)) + ".bin";
            try {
                std::filesystem::rename("last_input.bin", crashFile);
            } catch (const std::filesystem::filesystem_error& e) {
                std::cerr << "Failed to save crash input: " << e.what() << "\n";
            }

            throw 42;
        } else {
            return 1;
        }
    }
    std::ofstream(stateFile) << "1";
    ssize_t total_sent = 0;
    ssize_t to_send = inputBuf.size();
    const char* buffer = inputBuf.data();
    while (total_sent < to_send) {
        ssize_t sent = send(sock, buffer + total_sent, to_send - total_sent, 0);
        if (sent <= 0) {
            perror("send");
            close(sock);
            return 1;
        }
        total_sent += sent;
    }

    close(sock);

    std::ofstream out("last_input.bin", std::ios::binary);
    out.write(inputBuf.data(), inputBuf.size());
    out.close();


    return 0;
}
