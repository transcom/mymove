import React from 'react';
import { string, element, func, arrayOf, bool, shape } from 'prop-types';
import { Button } from '@trussworks/react-uswds';
import classnames from 'classnames';

import styles from './ReviewItems.module.scss';

const ReviewItems = ({ className, heading, renderAddButton, contents, emptyMessage }) => {
  return (
    <div className={classnames(styles.ReviewItems, className)}>
      <div className={styles.headingContainer}>
        <div className={styles.headingContent}>{heading}</div>
        {renderAddButton && <div className={styles.addButtonContainer}>{renderAddButton()}</div>}
      </div>
      <div className={styles.contentsContainer}>
        {(!contents || contents.length === 0) && (
          <div className={classnames({ [styles.subheadingWrapper]: !!renderAddButton }, 'display-flex', 'width-full')}>
            <span className={styles.emptyMessage}>{emptyMessage}</span>
          </div>
        )}
        {contents?.map(({ id, subheading, rows, onDelete, renderEditLink }) => {
          return (
            <div
              className={classnames({ [styles.subheadingWrapper]: !!renderAddButton }, 'display-flex', 'width-full')}
              key={id}
            >
              {subheading && <div className={styles.subheading}>{subheading}</div>}
              <dl>
                {rows.map(({ id: rowId, hideLabel, label, value }) => (
                  <div key={`${rowId}-${id}`}>
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
                    <span className={styles.actionSeparator}>|</span>
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
  className: string,
  heading: element.isRequired,
  renderAddButton: func,
  contents: arrayOf(
    shape({
      id: string.isRequired,
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
  ),
  emptyMessage: string,
};

ReviewItems.defaultProps = {
  className: '',
  renderAddButton: undefined,
  contents: undefined,
  emptyMessage: 'No items to display',
};

export default ReviewItems;
