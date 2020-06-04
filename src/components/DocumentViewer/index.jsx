import React from 'react';
import classNames from 'classnames/bind';

// import { ReactComponent as XLightIcon } from 'shared/icon/x-light.svg';
import { ReactComponent as DocMenu } from 'shared/icon/doc-menu.svg';
import fakeDoc from 'shared/images/fake-doc.png';
import styles from './index.module.scss';

const cx = classNames.bind(styles);

const sectionStyle = {
  backgroundImage: `url(${{ fakeDoc }})`,
};

const DocViewerMenu = () => (
  <div className={cx('doc-viewer--menu doc-viewer--menu--collapsed')}>
    <DocMenu />
    <div className={cx('thumbnail-container display-flex')}>
      <div className={cx('thumbnail-image')} style={sectionStyle} />
      <div className={cx('thumbnail-image')} style={sectionStyle} />
      <div className={cx('thumbnail-image')} style={sectionStyle} />
    </div>
  </div>
);

export default DocViewerMenu;
