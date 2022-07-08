import React from 'react';
import { string, element, func, arrayOf, bool, shape } from 'prop-types';
import { Button } from '@trussworks/react-uswds';
import classnames from 'classnames';

import styles from './ReviewItems.module.scss';

const ReviewItems = ({ heading, renderAddButton, contents }) => {
  return (
    <div className={styles.ReviewItems}>
      <div className={styles.headingContainer}>
        <div className={styles.headingContent}>{heading}</div>
        {renderAddButton && <div className={styles.addButtonContainer}>{renderAddButton()}</div>}
      </div>
      <div className={styles.contentsContainer}>
        {contents.map(({ subheading, rows, onDelete, renderEditLink }) => {
          return (
            <div className={classnames({ [styles.subheadingWrapper]: !!subheading }, 'display-flex', 'width-full')}>
              {subheading && <div className={styles.subheading}>{subheading}</div>}
              <dl>
                {rows.map(({ id, hideLabel, label, value }) => (
                  <div key={id}>
                    <dt className={classnames({ [styles.hiddenTerm]: hideLabel })}>{label}</dt>
                    <dd>{value}</dd>
                  </div>
                ))}
              </dl>
              <div className={styles.actionContainer}>
                {onDelete && (
                  <>
                    <Button type="button" unstyled onClick={onDelete}>
                      Delete
                    </Button>
                    <span className={styles.actionSeparater}>|</span>
                  </>
                )}
                {renderEditLink()}
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
};

ReviewItems.propTypes = {
  heading: element.isRequired,
  renderAddButton: func,
  contents: arrayOf(
    shape({
      subheading: element,
      rows: arrayOf(
        shape({
          id: string.isRequired,
          hideLabel: bool,
          label: string,
          value: string.isRequired,
        }),
      ).isRequired,
      onDelete: func,
      renderEditLink: func.isRequired,
    }),
  ).isRequired,
};

ReviewItems.defaultProps = {
  renderAddButton: undefined,
};

export default ReviewItems;
