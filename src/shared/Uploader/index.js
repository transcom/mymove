// eslint-disable-next-line no-unused-vars
import React, { Component } from 'react';

export class Uploader extends Component {
  constructor(props) {
    super(props);
    this.uploadFile = this.uploadFile.bind(this);
  }

  uploadFile() {
    console.log(this.fileInput.files[0]);
    // createDocument(this.fileInput.files[0])
  }
  render() {
    return (
      <div className="uploader">
        <input
          type="file"
          ref={input => {
            this.fileInput = input;
          }}
        />
        <button onClick={this.uploadFile}>Upload Now</button>
      </div>
    );
  }
}

export default Uploader;
