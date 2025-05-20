import React from 'react';
import classnames from 'classnames';

import styles from './FeedbackItems.module.scss';

import SectionWrapper from 'components/Shared/SectionWrapper/SectionWrapper';
import { DEFAULT_EMPTY_VALUE } from 'shared/constants';

const FeedbackItems = ({ documents, docType }) => {
  const formatSecondaryValue = (secondaryValue) => {
    return <span> (*{secondaryValue})</span>;
  };

  const formatDetails = (doc) => {
    return doc.map(({ label, value, secondaryValue }, i) => {
      if (value === DEFAULT_EMPTY_VALUE || !value) return null;
      return (
        <div className={styles.subheading} key={`${doc.key}-${i}`}>
          <span>{label}</span>
          <span>{value}</span>
          {secondaryValue && formatSecondaryValue(secondaryValue)}
        </div>
      );
    });
  };

  const formatHeader = (i) => {
    return (
      <div className={styles.subheading}>
        <h4 className="text-bold">
          {docType} {i + 1}
        </h4>
      </div>
    );
  };

  const formatFeedbackSet = (documentSet) => {
    return documentSet.map((doc, i) => {
      return (
        <SectionWrapper className={styles.headingWrapper} key={i}>
          <div className={styles.subheadingWrapper}>{formatHeader(i)}</div>
          <div>{formatDetails(doc)}</div>
        </SectionWrapper>
      );
    });
  };

  if (!documents) return null;

  const feedbackItems = formatFeedbackSet(documents);

  return (
    <div className={styles.FeedbackItems}>
      <div className={classnames(styles.contentsContainer)}>{feedbackItems}</div>
    </div>
  );
};

export default FeedbackItems;
