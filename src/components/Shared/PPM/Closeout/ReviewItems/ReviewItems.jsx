import React from 'react';
import { string, element, func, arrayOf, bool, shape, oneOfType, number, node } from 'prop-types';
import { Tag } from '@trussworks/react-uswds';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import classnames from 'classnames';

import styles from './ReviewItems.module.scss';

import { ButtonUsa as Button, destructiveOutlineButtonStyle } from 'shared/standardUI/Buttons/ButtonUsa';

const reviewWrapperStyle = styles['review-wrapper'];
const reviewPanelStyle = styles['review-panel'];
const reviewRowStyle = styles['review-row'];
const reviewDescriptionStyle = styles['review-description'];
const reviewTermStyle = styles['review-term'];
const reviewActionButtonContainer = styles['action-button-container'];
const reviewContentContainerStyle = styles['content-container'];
const ReviewItems = ({ className, heading, renderAddButton, contents, emptyMessage }) => {
  return (
    <div className={classnames(reviewWrapperStyle, className)}>
      <div className={styles.headingContainer}>
        <div className={styles.headingContent}>{heading}</div>
      </div>

      {(!contents || contents.length === 0) && (
        <div className={classnames({ [styles.emptyWrapper]: !!renderAddButton }, 'display-flex', 'width-full')}>
          <span className={styles.emptyMessage}>{emptyMessage}</span>
        </div>
      )}
      {renderAddButton && renderAddButton()}
      {contents?.map(({ id, isComplete, draftMessage, subheading, rows, onDelete, renderEditLink }, idx) => {
        return (
          <div className={reviewContentContainerStyle}>
            {isComplete === false && (
              <div className={styles.missingAlert}>
                <Tag className={classnames(styles.alertTag, 'usa-tag--alert')}>
                  <FontAwesomeIcon icon="exclamation" />
                </Tag>
                <span className="missingMessage">{draftMessage}</span>
              </div>
            )}
            {subheading && <div className={styles.subheading}>{subheading}</div>}
            <dl className={reviewPanelStyle}>
              {rows.map(({ id: rowId, hideLabel, label, value }) => (
                <div
                  key={`${rowId}-${id}`}
                  className={classnames({ [reviewRowStyle]: true, [styles[rowId]]: styles[rowId] })}
                >
                  <dt
                    className={classnames({ [reviewTermStyle]: true, [styles.hiddenTerm]: hideLabel })}
                    aria-hidden={hideLabel}
                  >
                    {label}
                  </dt>
                  <dd className={classnames({ [reviewDescriptionStyle]: true })}>{value}</dd>
                </div>
              ))}
            </dl>
            <div className={reviewActionButtonContainer}>
              {renderEditLink()}
              {onDelete && (
                <Button
                  className={destructiveOutlineButtonStyle}
                  data-testid={`weightMovedDelete-${idx + 1}`}
                  type="button"
                  onClick={onDelete}
                >
                  Delete
                </Button>
              )}
            </div>
          </div>
        );
      })}
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
          value: oneOfType([string, number, node]),
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
