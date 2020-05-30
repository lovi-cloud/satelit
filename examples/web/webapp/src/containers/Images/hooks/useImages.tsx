import { useState, useEffect } from "react";
import { SatelitClient } from "../../../clients/SatelitServiceClientPb";
import { Image, GetImagesRequest } from "../../../clients/satelit_pb";

export const useImages = (client: SatelitClient) => {
    const [isReload] = useState(0);
    const [images, setImages] = useState([] as Image[]);
    useEffect(() => reload(), [isReload]);

    const reload = () => {
        const request = new GetImagesRequest();
        client.getImages(request, {}, (err, res) => {
            if (err || res === null) {
                throw err;
            }
            const images: Image[] = res.getImagesList();
            setImages([...images]);
        });
    };

    return {
        images
    };
};