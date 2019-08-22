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
  imgEl = React.createRef();
  state = {
    rotation: 0,
    delta: 0,
    divHeight: null,
    divWidth: null,
    divXt: 0,
    ImgHeight: 'auto',
    ImgWidth: 'auto',
  };

  rotate = degrees => {
    let {delta} = this.state;
    delta = this.state.ImgWidth > this.state.ImgHeight ? - delta : delta;
    const radians = (this.state.rotation + degrees) / 180 * Math.PI;
    const times = Math.abs(Math.sin(radians));
    const [h, w] = [this.state.divWidth, this.state.divHeight];
    this.setState({
      divHeight: h,
      divWidth: w,
      Xt: delta * times,
      Yt: -delta * times,
      divXt: -2 * delta * times,
      rotation: (this.state.rotation + degrees) % 360,
    });
  };
  rotateLeft = () => {
    this.rotate(-90);
  };
  rotateRight = () => {
    this.rotate(90);
  };

  render() {
    const { divWidth, divHeight, rotation } = this.state;
    const adjDivWidth = divWidth ? divWidth + 10: divWidth;
    const adjDivHeight = divHeight ? divHeight + 80 : divWidth;
    let s = {};
    if (this.state.ImgWidth !== 'auto')
      s = (this.state.rotation === 90 || this.state.rotation === 180) ? {'maxWidth': `${this.state.ImgWidth}px`, 'maxHeight': 'none'} : {'maxHeight': `${this.state.ImgHeight}px`, 'maxWidth': 'none'};
    return (
      <div
        className="document-contents"
        style={{ padding: '5px', transform: `translateX(${this.state.divXt}px)`, width: adjDivWidth, height: adjDivHeight}}
      >
        <div style={{ marginBottom: '2em' }}>
          <RotationBar onLeftButtonClick={this.rotateLeft} onRightButtonClick={this.rotateRight} />
        </div>
        <div>
          <img
            className={styles[`rotate-${rotation}`]}
            style={{transform: `translate(${this.state.Xt}px, ${this.state.Yt}px) rotate(${rotation}deg)`, ...s}}
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
