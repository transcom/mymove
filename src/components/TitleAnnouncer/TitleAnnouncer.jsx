import React from 'react';

import styles from 'components/TitleAnnouncer/TitleAnnouncer.module.scss';

const TitleAnnouncer = () => <div id="title-announcer" aria-live="polite" className={styles.HiddenTitleAnnouncer} />;

export default TitleAnnouncer;
