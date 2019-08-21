import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { detectFirefox } from 'shared/utils';
import './index.css';
import styles from './DocumentContent.module.scss';
import { RotationBar } from 'shared/RotationBar/RotationBar';

const DocumentContent = props => {
  let { contentType, filename, url } = props;
  if (contentType === 'application/pdf') {
    return <PDFImage filename={filename} url={url} />;
  }
  return <NonPDFImage src={url} />;
};

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

export const PDFImage = props => {
  return (
    <div className="document-contents">
      <div className="page">
        {detectFirefox() ? (
          downloadOnlyView(props.filename, props.url)
        ) : (
          <object className={styles.pdf} data={props.url} type="application/pdf" alt="document upload">
            {downloadOnlyView(props.filename, props.url)}
          </object>
        )}
      </div>
    </div>
  );
};

PDFImage.propTypes = {
  filename: PropTypes.string.isRequired,
  url: PropTypes.string.isRequired,
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
        <div>
          <img className={styles[`rotate-${this.state.rotation}`]} src={this.props.src} alt="document upload" />
        </div>
      </div>
    );
  }
}

NonPDFImage.propTypes = {
  src: PropTypes.string.isRequired,
};

export default DocumentContent;
