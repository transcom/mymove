import React, { Component } from 'react';
import PropTypes from 'prop-types';

import './index.css';

export class DocumentContent extends Component {
  imgEl = React.createRef();

  state = {
    orientation: this.props.orientation,
    imgHeight: 0,
  };

  adjustContainerHeight() {
    const imgHeight = this.imgEl.current.getBoundingClientRect().height;
    this.setState({
      imgHeight: imgHeight,
    });
  }

  handleRotate(direction) {
    if (direction === 'right') {
      console.log(this.state.orientation, this.state.orientation - 90);
      this.setState({ orientation: this.state.orientation - 90 });
      this.props.rotate(this.props.uploadId, this.state.orientation - 90);
    } else {
      console.log(this.state.orientation, this.state.orientation + 90);
      this.setState({ orientation: this.state.orientation + 90 });
      this.props.rotate(this.props.uploadId, this.state.orientation + 90);
    }
  }

  render() {
    return (
      <div
        className="page"
        style={{
          minHeight: this.state.imgHeight + 50,
          display: 'flex',
          flexDirection: 'column',
          justifyContent: 'space-between',
        }}
      >
        {this.props.contentType === 'application/pdf' ? (
          <div className="pdf-placeholder">
            {this.props.filename && <span className="filename">{this.props.filename}</span>}
            This PDF can be{' '}
            <a target="_blank" href={this.props.url}>
              viewed here
            </a>
            .
          </div>
        ) : (
          <div style={{ marginTop: this.state.imgHeight / 5, marginBottom: this.state.imgHeight / 5 }}>
            <img
              src={this.props.url}
              ref={this.imgEl}
              style={{ transform: `rotate(${this.state.orientation}deg)` }}
              onLoad={this.adjustContainerHeight.bind(this)}
              alt="document upload"
            />
          </div>
        )}

        <div>
          <button onClick={this.handleRotate.bind(this, 'left')}>rotate left</button>
          <button onClick={this.handleRotate.bind(this, 'right')}>rotate right</button>
        </div>
      </div>
    );
  }
}

DocumentContent.propTypes = {
  contentType: PropTypes.string.isRequired,
  filename: PropTypes.string.isRequired,
  url: PropTypes.string.isRequired,
  uploadId: PropTypes.string.isRequired,
  orientation: PropTypes.number,
  rotate: PropTypes.func.isRequired,
};

export default DocumentContent;
