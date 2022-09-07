import React from 'react';
import PropTypes from 'prop-types';
import { Grid, Accordion, Checkbox } from '@trussworks/react-uswds';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import styles from './ViolationsAccordion.module.scss';

import { PWSViolationShape } from 'types/pwsViolation';

/**
 * A component that renders an accordion for a single violation category. An expandable section of the accordion is rendered for each subcategory.
 */
const ViolationsAccordion = ({ violations, onChange, selected }) => {
  const [expandedViolations, setExpandedViolations] = React.useState([]);
  const subCategories = [...new Set(violations.map((item) => item.subCategory))];

  const toggleDetailExpand = (violationId) => {
    if (expandedViolations.includes(violationId)) {
      setExpandedViolations(expandedViolations.filter((id) => id !== violationId));
    } else {
      setExpandedViolations([...expandedViolations, violationId]);
    }
  };

  const getContentForItem = (subCategory) => {
    const subCategoryViolations = violations.filter((violation) => violation.subCategory === subCategory);
    const items = subCategoryViolations.map((violation) => (
      <div key={`${violation.id}-accordion-option`} className={styles.accordionOption}>
        <div className={styles.flex}>
          <Checkbox
            id={`${violation.id}-checkbox`}
            name={`${violation.paragraphNumber} ${violation.title}`}
            className={styles.checkbox}
            aria-labelledby={`${violation.id}-checkbox-label`}
            onChange={() => onChange(violation.id)}
            checked={!!(selected && selected.includes(violation.id))}
          />

          {/* Checkbox label */}
          <div className={styles.grow} id={`${violation.id}-checkbox-label`}>
            <h5>{`${violation.paragraphNumber} ${violation.title}`}</h5>
            <small>{violation.requirementSummary}</small>
          </div>

          {/* Expand Requirements Statement Toggle Button */}
          {expandedViolations.includes(violation.id) ? (
            <FontAwesomeIcon
              icon="chevron-down"
              className={styles.detailIcon}
              role="button"
              onClick={() => {
                toggleDetailExpand(violation.id);
              }}
              fontSize="20px"
              data-testid="collapse-icon"
              aria-label="Collapse Requirements Statement"
            />
          ) : (
            <FontAwesomeIcon
              icon="chevron-up"
              className={styles.detailIcon}
              role="button"
              onClick={() => {
                toggleDetailExpand(violation.id);
              }}
              fontSize="20px"
              data-testid="expand-icon"
              aria-label="Expand Requirements Statement"
            />
          )}
        </div>

        {/* Expandable Requirements Statement */}
        {expandedViolations.includes(violation.id) && (
          <p className={styles.requirementStatement}>
            <small>{violation.requirementStatement}</small>
          </p>
        )}
      </div>
    ));

    return items;
  };

  const getAccordionItems = () => {
    const items = [];
    subCategories.forEach((subCategory) => {
      items.push({
        title: subCategory,
        content: getContentForItem(subCategory),
        expanded: false,
        id: `${subCategory}-violation`,
        headingLevel: 'h4',
      });
    });

    return items;
  };

  const { category } = violations[0]; // All should be the same violaiton category
  return (
    <>
      <Grid row key={`${category}-accordion-category`}>
        <Grid col>
          <h3>{category}</h3>
        </Grid>
      </Grid>
      <div>
        <Accordion items={getAccordionItems()} multiselectable bordered className={styles.accordion} />
      </div>
    </>
  );
};

ViolationsAccordion.propTypes = {
  violations: PropTypes.arrayOf(PWSViolationShape),
};

ViolationsAccordion.defaultProps = {
  violations: [],
};

export default ViolationsAccordion;
