/**
 * @fileoverview gRPC-Web generated client stub for satelit
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!
/* eslint-disable */


import * as grpcWeb from 'grpc-web';

import {
  AddVolumeRequest,
  AddVolumeResponse,
  AttachVolumeRequest,
  AttachVolumeResponse,
  DeleteImageRequest,
  DeleteImageResponse,
  GetImagesRequest,
  GetImagesResponse,
  GetVolumesRequest,
  GetVolumesResponse,
  UploadImageRequest,
  UploadImageResponse} from './satelit_pb';

export class SatelitClient {
  client_: grpcWeb.AbstractClientBase;
  hostname_: string;
  credentials_: null | { [index: string]: string; };
  options_: null | { [index: string]: string; };

  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: string; }) {
    if (!options) options = {};
    if (!credentials) credentials = {};
    options['format'] = 'text';

    this.client_ = new grpcWeb.GrpcWebClientBase(options);
    this.hostname_ = hostname;
    this.credentials_ = credentials;
    this.options_ = options;
  }

  methodInfoGetVolumes = new grpcWeb.AbstractClientBase.MethodInfo(
    GetVolumesResponse,
    (request: GetVolumesRequest) => {
      return request.serializeBinary();
    },
    GetVolumesResponse.deserializeBinary
  );

  getVolumes(
    request: GetVolumesRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.Error,
               response: GetVolumesResponse) => void) {
    return this.client_.rpcCall(
      this.hostname_ +
        '/satelit.Satelit/GetVolumes',
      request,
      metadata || {},
      this.methodInfoGetVolumes,
      callback);
  }

  methodInfoAddVolume = new grpcWeb.AbstractClientBase.MethodInfo(
    AddVolumeResponse,
    (request: AddVolumeRequest) => {
      return request.serializeBinary();
    },
    AddVolumeResponse.deserializeBinary
  );

  addVolume(
    request: AddVolumeRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.Error,
               response: AddVolumeResponse) => void) {
    return this.client_.rpcCall(
      this.hostname_ +
        '/satelit.Satelit/AddVolume',
      request,
      metadata || {},
      this.methodInfoAddVolume,
      callback);
  }

  methodInfoAttachVolume = new grpcWeb.AbstractClientBase.MethodInfo(
    AttachVolumeResponse,
    (request: AttachVolumeRequest) => {
      return request.serializeBinary();
    },
    AttachVolumeResponse.deserializeBinary
  );

  attachVolume(
    request: AttachVolumeRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.Error,
               response: AttachVolumeResponse) => void) {
    return this.client_.rpcCall(
      this.hostname_ +
        '/satelit.Satelit/AttachVolume',
      request,
      metadata || {},
      this.methodInfoAttachVolume,
      callback);
  }

  methodInfoGetImages = new grpcWeb.AbstractClientBase.MethodInfo(
    GetImagesResponse,
    (request: GetImagesRequest) => {
      return request.serializeBinary();
    },
    GetImagesResponse.deserializeBinary
  );

  getImages(
    request: GetImagesRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.Error,
               response: GetImagesResponse) => void) {
    return this.client_.rpcCall(
      this.hostname_ +
        '/satelit.Satelit/GetImages',
      request,
      metadata || {},
      this.methodInfoGetImages,
      callback);
  }

  methodInfoDeleteImage = new grpcWeb.AbstractClientBase.MethodInfo(
    DeleteImageResponse,
    (request: DeleteImageRequest) => {
      return request.serializeBinary();
    },
    DeleteImageResponse.deserializeBinary
  );

  deleteImage(
    request: DeleteImageRequest,
    metadata: grpcWeb.Metadata | null,
    callback: (err: grpcWeb.Error,
               response: DeleteImageResponse) => void) {
    return this.client_.rpcCall(
      this.hostname_ +
        '/satelit.Satelit/DeleteImage',
      request,
      metadata || {},
      this.methodInfoDeleteImage,
      callback);
  }

}

