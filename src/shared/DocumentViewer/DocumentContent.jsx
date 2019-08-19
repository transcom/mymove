import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { detectFirefox } from 'shared/utils';
import './index.css';
import styles from './DocumentContent.module.scss';
import { RotationBar } from 'shared/RotationBar/RotationBar';

class DocumentContent extends Component {
  render() {
    let { contentType, filename, url } = this.props;
    if (contentType === 'application/pdf') {
      return <PDFImage filename={filename} url={url} />;
    }
    return <NonPDFImage src={url} />;
  }
}

DocumentContent.propTypes = {
  contentType: PropTypes.string.isRequired,
  filename: PropTypes.string.isRequired,
  url: PropTypes.string.isRequired,
};

const downloadOnlyView = (filename, url) => (
  <div className="pdf-placeholder">
    {filename && <span className="filename">{filename}</span>}
    This PDF can be{' '}
    <a target="_blank" rel="noopener noreferrer" href={url}>
      viewed here
    </a>
    .
  </div>
);

export class PDFImage extends Component {
  render() {
    return (
      <div className="document-contents">
        <div className="page">
          {detectFirefox() ? (
            downloadOnlyView(this.props.filename, this.props.url)
          ) : (
            <object className={styles.pdf} data={this.props.url} type="application/pdf" alt="document upload">
              {downloadOnlyView(this.props.filename, this.props.url)}
            </object>
          )}
        </div>
      </div>
    );
  }
}

PDFImage.propTypes = {
  filename: PropTypes.any,
  url: PropTypes.any,
};

export class NonPDFImage extends Component {
  state = {
    rotation: 0,
  };

  rotate = degrees => {
    this.setState({
      rotation: (360 + this.state.rotation + degrees) % 360,
    });
  };

  rotateLeft = () => {
    this.rotate(-90);
  };
  rotateRight = () => {
    this.rotate(90);
  };

  render() {
    return (
      <div className="document-contents" style={{ padding: '1em' }}>
        <div style={{ marginBottom: '2em' }}>
          <RotationBar onLeftButtonClick={this.rotateLeft} onRightButtonClick={this.rotateRight} />
        </div>
        <div className="non-pdf-img-container">
          <img className={`non-pdf-img rotate-${this.state.rotation}`} src={this.props.src} alt="document upload" />
        </div>
      </div>
    );
  }
}

NonPDFImage.propTypes = {
  onClick: PropTypes.func,
  src: PropTypes.any,
};

export default DocumentContent;
