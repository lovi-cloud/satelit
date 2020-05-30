import * as jspb from "google-protobuf"

export class Volume extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getAttached(): boolean;
  setAttached(value: boolean): void;

  getHostname(): string;
  setHostname(value: string): void;

  getCapacityByte(): number;
  setCapacityByte(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Volume.AsObject;
  static toObject(includeInstance: boolean, msg: Volume): Volume.AsObject;
  static serializeBinaryToWriter(message: Volume, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Volume;
  static deserializeBinaryFromReader(message: Volume, reader: jspb.BinaryReader): Volume;
}

export namespace Volume {
  export type AsObject = {
    id: string,
    attached: boolean,
    hostname: string,
    capacityByte: number,
  }
}

export class Image extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getName(): string;
  setName(value: string): void;

  getVolumeId(): string;
  setVolumeId(value: string): void;

  getDescription(): string;
  setDescription(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Image.AsObject;
  static toObject(includeInstance: boolean, msg: Image): Image.AsObject;
  static serializeBinaryToWriter(message: Image, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Image;
  static deserializeBinaryFromReader(message: Image, reader: jspb.BinaryReader): Image;
}

export namespace Image {
  export type AsObject = {
    id: string,
    name: string,
    volumeId: string,
    description: string,
  }
}

export class GetVolumesRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetVolumesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetVolumesRequest): GetVolumesRequest.AsObject;
  static serializeBinaryToWriter(message: GetVolumesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetVolumesRequest;
  static deserializeBinaryFromReader(message: GetVolumesRequest, reader: jspb.BinaryReader): GetVolumesRequest;
}

export namespace GetVolumesRequest {
  export type AsObject = {
  }
}

export class GetVolumesResponse extends jspb.Message {
  getVolumesList(): Array<Volume>;
  setVolumesList(value: Array<Volume>): void;
  clearVolumesList(): void;
  addVolumes(value?: Volume, index?: number): Volume;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetVolumesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetVolumesResponse): GetVolumesResponse.AsObject;
  static serializeBinaryToWriter(message: GetVolumesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetVolumesResponse;
  static deserializeBinaryFromReader(message: GetVolumesResponse, reader: jspb.BinaryReader): GetVolumesResponse;
}

export namespace GetVolumesResponse {
  export type AsObject = {
    volumesList: Array<Volume.AsObject>,
  }
}

export class AddVolumeRequest extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  getCapacityByte(): number;
  setCapacityByte(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddVolumeRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AddVolumeRequest): AddVolumeRequest.AsObject;
  static serializeBinaryToWriter(message: AddVolumeRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddVolumeRequest;
  static deserializeBinaryFromReader(message: AddVolumeRequest, reader: jspb.BinaryReader): AddVolumeRequest;
}

export namespace AddVolumeRequest {
  export type AsObject = {
    name: string,
    capacityByte: number,
  }
}

export class AddVolumeResponse extends jspb.Message {
  getVolume(): Volume | undefined;
  setVolume(value?: Volume): void;
  hasVolume(): boolean;
  clearVolume(): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AddVolumeResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AddVolumeResponse): AddVolumeResponse.AsObject;
  static serializeBinaryToWriter(message: AddVolumeResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AddVolumeResponse;
  static deserializeBinaryFromReader(message: AddVolumeResponse, reader: jspb.BinaryReader): AddVolumeResponse;
}

export namespace AddVolumeResponse {
  export type AsObject = {
    volume?: Volume.AsObject,
  }
}

export class AttachVolumeRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  getHostname(): string;
  setHostname(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AttachVolumeRequest.AsObject;
  static toObject(includeInstance: boolean, msg: AttachVolumeRequest): AttachVolumeRequest.AsObject;
  static serializeBinaryToWriter(message: AttachVolumeRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AttachVolumeRequest;
  static deserializeBinaryFromReader(message: AttachVolumeRequest, reader: jspb.BinaryReader): AttachVolumeRequest;
}

export namespace AttachVolumeRequest {
  export type AsObject = {
    id: string,
    hostname: string,
  }
}

export class AttachVolumeResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): AttachVolumeResponse.AsObject;
  static toObject(includeInstance: boolean, msg: AttachVolumeResponse): AttachVolumeResponse.AsObject;
  static serializeBinaryToWriter(message: AttachVolumeResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): AttachVolumeResponse;
  static deserializeBinaryFromReader(message: AttachVolumeResponse, reader: jspb.BinaryReader): AttachVolumeResponse;
}

export namespace AttachVolumeResponse {
  export type AsObject = {
  }
}

export class DeleteVolumeRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteVolumeRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteVolumeRequest): DeleteVolumeRequest.AsObject;
  static serializeBinaryToWriter(message: DeleteVolumeRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteVolumeRequest;
  static deserializeBinaryFromReader(message: DeleteVolumeRequest, reader: jspb.BinaryReader): DeleteVolumeRequest;
}

export namespace DeleteVolumeRequest {
  export type AsObject = {
  }
}

export class DeleteVolumeResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteVolumeResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteVolumeResponse): DeleteVolumeResponse.AsObject;
  static serializeBinaryToWriter(message: DeleteVolumeResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteVolumeResponse;
  static deserializeBinaryFromReader(message: DeleteVolumeResponse, reader: jspb.BinaryReader): DeleteVolumeResponse;
}

export namespace DeleteVolumeResponse {
  export type AsObject = {
  }
}

export class GetImagesRequest extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetImagesRequest.AsObject;
  static toObject(includeInstance: boolean, msg: GetImagesRequest): GetImagesRequest.AsObject;
  static serializeBinaryToWriter(message: GetImagesRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetImagesRequest;
  static deserializeBinaryFromReader(message: GetImagesRequest, reader: jspb.BinaryReader): GetImagesRequest;
}

export namespace GetImagesRequest {
  export type AsObject = {
  }
}

export class GetImagesResponse extends jspb.Message {
  getImagesList(): Array<Image>;
  setImagesList(value: Array<Image>): void;
  clearImagesList(): void;
  addImages(value?: Image, index?: number): Image;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): GetImagesResponse.AsObject;
  static toObject(includeInstance: boolean, msg: GetImagesResponse): GetImagesResponse.AsObject;
  static serializeBinaryToWriter(message: GetImagesResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): GetImagesResponse;
  static deserializeBinaryFromReader(message: GetImagesResponse, reader: jspb.BinaryReader): GetImagesResponse;
}

export namespace GetImagesResponse {
  export type AsObject = {
    imagesList: Array<Image.AsObject>,
  }
}

export class UploadImageRequest extends jspb.Message {
  getMeta(): UploadImageRequestMeta | undefined;
  setMeta(value?: UploadImageRequestMeta): void;
  hasMeta(): boolean;
  clearMeta(): void;

  getChunk(): UploadImageRequestChunk | undefined;
  setChunk(value?: UploadImageRequestChunk): void;
  hasChunk(): boolean;
  clearChunk(): void;

  getValueCase(): UploadImageRequest.ValueCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UploadImageRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UploadImageRequest): UploadImageRequest.AsObject;
  static serializeBinaryToWriter(message: UploadImageRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UploadImageRequest;
  static deserializeBinaryFromReader(message: UploadImageRequest, reader: jspb.BinaryReader): UploadImageRequest;
}

export namespace UploadImageRequest {
  export type AsObject = {
    meta?: UploadImageRequestMeta.AsObject,
    chunk?: UploadImageRequestChunk.AsObject,
  }

  export enum ValueCase { 
    VALUE_NOT_SET = 0,
    META = 1,
    CHUNK = 2,
  }
}

export class UploadImageRequestMeta extends jspb.Message {
  getName(): string;
  setName(value: string): void;

  getDescription(): string;
  setDescription(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UploadImageRequestMeta.AsObject;
  static toObject(includeInstance: boolean, msg: UploadImageRequestMeta): UploadImageRequestMeta.AsObject;
  static serializeBinaryToWriter(message: UploadImageRequestMeta, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UploadImageRequestMeta;
  static deserializeBinaryFromReader(message: UploadImageRequestMeta, reader: jspb.BinaryReader): UploadImageRequestMeta;
}

export namespace UploadImageRequestMeta {
  export type AsObject = {
    name: string,
    description: string,
  }
}

export class UploadImageRequestChunk extends jspb.Message {
  getData(): Uint8Array | string;
  getData_asU8(): Uint8Array;
  getData_asB64(): string;
  setData(value: Uint8Array | string): void;

  getPosition(): number;
  setPosition(value: number): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UploadImageRequestChunk.AsObject;
  static toObject(includeInstance: boolean, msg: UploadImageRequestChunk): UploadImageRequestChunk.AsObject;
  static serializeBinaryToWriter(message: UploadImageRequestChunk, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UploadImageRequestChunk;
  static deserializeBinaryFromReader(message: UploadImageRequestChunk, reader: jspb.BinaryReader): UploadImageRequestChunk;
}

export namespace UploadImageRequestChunk {
  export type AsObject = {
    data: Uint8Array | string,
    position: number,
  }
}

export class UploadImageResponse extends jspb.Message {
  getImage(): Image | undefined;
  setImage(value?: Image): void;
  hasImage(): boolean;
  clearImage(): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UploadImageResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UploadImageResponse): UploadImageResponse.AsObject;
  static serializeBinaryToWriter(message: UploadImageResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UploadImageResponse;
  static deserializeBinaryFromReader(message: UploadImageResponse, reader: jspb.BinaryReader): UploadImageResponse;
}

export namespace UploadImageResponse {
  export type AsObject = {
    image?: Image.AsObject,
  }
}

export class DeleteImageRequest extends jspb.Message {
  getId(): string;
  setId(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteImageRequest.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteImageRequest): DeleteImageRequest.AsObject;
  static serializeBinaryToWriter(message: DeleteImageRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteImageRequest;
  static deserializeBinaryFromReader(message: DeleteImageRequest, reader: jspb.BinaryReader): DeleteImageRequest;
}

export namespace DeleteImageRequest {
  export type AsObject = {
    id: string,
  }
}

export class DeleteImageResponse extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): DeleteImageResponse.AsObject;
  static toObject(includeInstance: boolean, msg: DeleteImageResponse): DeleteImageResponse.AsObject;
  static serializeBinaryToWriter(message: DeleteImageResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): DeleteImageResponse;
  static deserializeBinaryFromReader(message: DeleteImageResponse, reader: jspb.BinaryReader): DeleteImageResponse;
}

export namespace DeleteImageResponse {
  export type AsObject = {
  }
}

