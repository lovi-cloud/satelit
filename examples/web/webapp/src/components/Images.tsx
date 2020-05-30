import React from "react";
import { Image } from "../clients/satelit_pb"

type Props = {
    images: Image[];
};

export const Images: React.FC<Props> = ({ images }: Props) => {
    return (
        <div>
            {images.map(image => (
                <div key={image.getId()}>{image.getName()}</div>
            ))}
        </div>
    )
}