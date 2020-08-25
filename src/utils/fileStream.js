export const SLICE_SIZE = 32 * 1024 * 1024;

export function defer() {
  let resolve, reject
  const p = new Promise(
    (res, rej) => ([resolve, reject] = [res, rej]),
  );
  return Object.assign(p, {
    resolve,
    reject,
  });
}

export class FileStream {
  constructor(file) {
    this.file = file;
    this.reader = new FileReader()
    this.bytesRead = 0
    this.pendingRead = undefined
    this.reader.onloadend = () => this.onLoad()
  }

  onLoad() {
    const res = this.reader.result;
    const pendingRead = this.pendingRead;
    this.pendingRead = undefined;
    if (this.reader.error) {
      pendingRead.reject(this.reader.error);
      return;
    }
    this.bytesRead += res.byteLength;
    pendingRead.resolve({
      data: new Uint8Array(res),
      eof: this.bytesRead >= this.file.size,
      bytesRead: this.bytesRead,
      bytesTotal: this.file.size,
      length: this.file.size - this.bytesRead,
    });
  }

  readChunk() {
    const sliceEnd = Math.min(this.bytesRead + SLICE_SIZE, this.file.size);
    const slice = this.file.slice(this.bytesRead, sliceEnd);
    this.pendingRead = defer();
    this.reader.readAsArrayBuffer(slice);
    return this.pendingRead;
  }
}
