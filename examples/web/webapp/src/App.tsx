import React from "react";
import { ImageContainer } from "./containers/Images";
import { gRPCClients } from "./gRPCClients";

export const App = () => {
    return (
        <div>
            <ImageContainer clients={gRPCClients} />
        </div>
    );
};