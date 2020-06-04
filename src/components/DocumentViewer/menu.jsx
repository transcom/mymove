import React from 'react';
import classNames from 'classnames/bind';

import { ReactComponent as XLightIcon } from 'shared/icon/x-light.svg';
import { ReactComponent as DocMenu } from 'shared/icon/doc-menu.svg';
import styles from './index.module.scss';

const cx = classNames.bind(styles);

const DocViewerMenu = () => (
  <div className={cx('doc-viewer--menu doc-viewer--menu--collapsed')}>
    <DocMenu />
    <XLightIcon />
    <div className={cx('thumbnail-container')}>test</div>
  </div>
);

export default DocViewerMenu;
