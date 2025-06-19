import React, { Component } from 'react';
import PropTypes from 'prop-types';

import { detectFirefox } from 'shared/utils';

import './index.css';
import styles from './DocumentContent.module.scss';

import { RotationBar } from 'shared/RotationBar/RotationBar';
import Alert from 'shared/Alert';
import { UPLOAD_SCAN_STATUS } from 'shared/constants';

const DocumentContent = (props) => {
  const { contentType, filename, url, status } = props;

  if (status === UPLOAD_SCAN_STATUS.THREATS_FOUND || status === UPLOAD_SCAN_STATUS.LEGACY_INFECTED) {
    return (
      <Alert type="error" className="usa-width-one-whole" heading="Ask for a new file">
        Our antivirus software flagged this file as a security risk. Contact the service member. Ask them to upload a
        photo of the original document instead.
      </Alert>
    );
  }
  if (status === UPLOAD_SCAN_STATUS.PROCESSING) {
    return (
      <Alert type="info" className="usa-width-one-whole" heading="Your file is being scanned for viruses">
        It will be available within a few minutes.
      </Alert>
    );
  }
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
    <a target="_blank" rel="noopener noreferrer" href={url} className="usa-link">
      viewed here
    </a>
    .
  </div>
);

export const PDFImage = (props) => {
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
  imgEl = React.createRef();

  state = {
    rotation: 0,
    delta: 0,
    divHeight: null,
    divWidth: null,
    divXt: 0,
    ImgHeight: 'auto',
    ImgWidth: 'auto',
    maxHeight: 'auto',
    maxWidth: 'auto',
  };

  transformImage = (degrees) => {
    const rotation = this.rotate(degrees);
    const sign = this.translateSign(rotation);
    const imageTranslationToApply = this.translateImage(sign);
    const containerTranslationToApply = this.translateContainer(sign);
    const imageMaxHeightWidth = this.calcImgMaxHeight();
    this.setState({
      rotation,
      sign,
      ...imageTranslationToApply,
      ...containerTranslationToApply,
      ...imageMaxHeightWidth,
    });
  };

  rotate = (degrees) => {
    return (this.state.rotation + degrees) % 360;
  };

  rotateLeft = () => {
    this.transformImage(-90);
  };

  rotateRight = () => {
    this.transformImage(90);
  };

  // delta between height + width. used for calculating the image + container translations
  get delta() {
    const { delta } = this.state;
    return this.state.ImgWidth > this.state.ImgHeight ? -delta : delta;
  }

  translateImage = (sign) => {
    return {
      Xt: this.delta * sign,
      Yt: -this.delta * sign,
    };
  };

  translateContainer = (sign) => {
    const { divHeight, divWidth } = this.state;
    // switch divWidth and divHeight
    return {
      divHeight: divWidth,
      divWidth: divHeight,
      divXt: -2 * this.delta * sign,
    };
  };

  // rotation dependent sign to be used with coordinate translations
  translateSign = (rotation) => {
    const radians = (rotation / 180) * Math.PI;
    return Math.abs(Math.sin(radians));
  };

  calcImgMaxHeight = () => {
    const { ImgWidth, ImgHeight, rotation } = this.state;
    const { maxHeight, maxWidth } = this.state;
    let maxHeightWidth = { maxHeight, maxWidth };
    if (ImgWidth !== 'auto') {
      maxHeightWidth =
        rotation === 90 || rotation === 180
          ? { maxWidth: ImgWidth, maxHeight: 'none' }
          : { maxWidth: 'none', maxHeight: ImgHeight };
    }
    return maxHeightWidth;
  };

  formatMaxes = (max) => {
    return max === 'none' ? 'none' : `${max}px`;
  };

  render() {
    const { divWidth, divHeight, rotation } = this.state;
    const adjDivWidth = divWidth ? divWidth + 10 : divWidth;
    const adjDivHeight = divHeight ? divHeight + 80 : divWidth;
    return (
      <div
        className="document-contents"
        style={{
          padding: '5px',
          transform: `translateX(${this.state.divXt}px)`,
          width: adjDivWidth,
          height: adjDivHeight,
        }}
      >
        <div style={{ marginBottom: '2em' }}>
          <RotationBar onLeftButtonClick={this.rotateLeft} onRightButtonClick={this.rotateRight} />
        </div>
        <div>
          <img
            style={{
              transform: `translate(${this.state.Xt}px, ${this.state.Yt}px) rotate(${rotation}deg)`,
              maxHeight: this.formatMaxes(this.state.maxHeight),
              maxWidth: this.formatMaxes(this.state.maxWidth),
            }}
            src={this.props.src}
            alt="document upload"
            ref={this.imgEl}
            onLoad={() =>
              this.setState({
                divHeight: this.imgEl.current.height,
                divWidth: this.imgEl.current.width,
                ImgHeight: this.imgEl.current.height,
                ImgWidth: this.imgEl.current.width,
                delta: Math.abs(this.imgEl.current.height - this.imgEl.current.width) / 2,
              })
            }
          />
        </div>
      </div>
    );
  }
}

NonPDFImage.propTypes = {
  src: PropTypes.string.isRequired,
};

export default DocumentContent;
