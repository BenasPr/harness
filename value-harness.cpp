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
#include <cstring>
#include "tuning.pb.h"

int main(int argc, char* argv[]) {
    std::vector<char> inputBuf((std::istreambuf_iterator<char>(std::cin)),
                                std::istreambuf_iterator<char>());

    std::vector<double> values(7);
    std::memcpy(values.data(), inputBuf.data(), 56);

    std::vector<std::string> keys = {
        "speed",
        "kp",
        "kd",
        "ki",
        "servo-trim",
        "threshhold-value",
        "servo-scaler"
    };

    protobuf_msgs::TuningState tuning;
    for (size_t i = 0; i < 7; ++i) {
        auto* param = tuning.add_dynamicparameters();
        auto* numberParam = new protobuf_msgs::TuningState::Parameter::NumberParameter();
        numberParam->set_key(keys[i]);
        numberParam->set_value(values[i]);
        param->set_allocated_number(numberParam);
    }

    std::string serialized;
    if (!tuning.SerializeToString(&serialized)) {
        return 1;
    }

    std::string stateFile = "pipeline_status.flag";
    bool wasConnectedBefore = false;
    if (std::ifstream in{stateFile}) {
        char c;
        in >> c;
        wasConnectedBefore = (c == '1');
    }
 

    int sock = socket(AF_INET, SOCK_STREAM, 0);
    if (sock < 0) {
        return 1;
    }

    sockaddr_in serv_addr{};
    serv_addr.sin_family = AF_INET;
    serv_addr.sin_port = htons(9000);
    if (inet_pton(AF_INET, "192.168.0.146", &serv_addr.sin_addr) <= 0) {
        close(sock);
        return 1;
    }

    if (connect(sock, (sockaddr*)&serv_addr, sizeof(serv_addr)) < 0) {
        std::ofstream(stateFile) << "0";
        if (wasConnectedBefore) {
            std::filesystem::create_directory("crashes");

            std::string crashFile = "crashes/crash_" + std::to_string(std::time(nullptr)) + ".bin";
            std::string scpCommand = "sshpass -p 'debix' scp debix@192.168.0.146:/tmp/last_input.bin " + crashFile;
            std::system(scpCommand.c_str());
        } else {
            system("./start_pipeline.sh");
        }return 1;
    }

    std::ofstream(stateFile) << "1";

    ssize_t total_sent = 0;
    ssize_t to_send = serialized.size();
    const char* buffer = serialized.data();
    while (total_sent < to_send) {
        ssize_t sent = send(sock, buffer + total_sent, to_send - total_sent, 0);
        if (sent <= 0) {
            close(sock);
            return 1;
        }
        total_sent += sent;
    }

    close(sock);

    // std::ofstream out("last_input.bin", std::ios::binary);
    // out.write(inputBuf.data(), inputBuf.size());
    // out.close();

    return 0;
}
