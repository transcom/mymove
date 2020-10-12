import React from 'react';
import classNames from 'classnames/bind';

import styles from './SectionWrapper.module.scss';

const cx = classNames.bind(styles);

const SectionWrapper = ({ children }) => <div className={cx('container-co')}>{children}</div>;

export { SectionWrapper as SectionWrapperComponent };