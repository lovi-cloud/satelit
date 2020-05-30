/* eslint-disable */
import React from "react";
import { Images } from "../../components/Images";
import { GRPCClients } from "../../gRPCClients";
import { useImages } from "./hooks/useImages";

type Props = {
    clients: GRPCClients;
}

export const ImageContainer: React.FC<Props> = ({ clients }: Props) => {
    const imageClient = clients.satelitClient;
    const imageState = useImages(imageClient);

    return (
        <div>
            <Images {...imageState}/>
        </div>
    )
}