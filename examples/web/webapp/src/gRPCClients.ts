import { SatelitClient } from "./clients/SatelitServiceClientPb";

export type GRPCClients = {
    satelitClient: SatelitClient;
}

export const gRPCClients = {
    satelitClient: new SatelitClient("http://localhost:8080") // running Envoy proxy on localhost:8080
}