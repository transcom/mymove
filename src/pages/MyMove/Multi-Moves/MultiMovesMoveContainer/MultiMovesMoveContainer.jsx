import React, { useState } from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import classnames from 'classnames';
import { Button } from '@trussworks/react-uswds';

import MultiMovesMoveInfoList from '../MultiMovesMoveInfoList/MultiMovesMoveInfoList';

import styles from './MultiMovesMoveContainer.module.scss';

import ButtonDropdown from 'components/ButtonDropdown/ButtonDropdown';

const MultiMovesMoveContainer = ({ move }) => {
  const [expandedMoves, setExpandedMoves] = useState({});

  const handleExpandClick = (index) => {
    setExpandedMoves((prev) => ({
      ...prev,
      [index]: !prev[index],
    }));
  };

  const handleButtonDropdownChange = (e) => {
    const selectedOption = e.target.value;
    console.log(selectedOption);
  };

  const moveList = move.map((m, index) => (
    <React.Fragment key={index}>
      <div className={styles.moveContainer}>
        <div className={styles.heading} key={index}>
          <h3>#{m.moveCode}</h3>
          {m.status !== 'APPROVED' ? (
            <Button className={styles.moveHeaderBtn}>Go to Move</Button>
          ) : (
            <ButtonDropdown divClassName={styles.moveHeaderBtn} onChange={handleButtonDropdownChange}>
              <option value="">Download</option>
              <option>PCS Orders</option>
              <option>PPM Packet</option>
            </ButtonDropdown>
          )}
          <FontAwesomeIcon
            className={styles.icon}
            icon={classnames({
              'chevron-up': expandedMoves[index],
              'chevron-down': !expandedMoves[index],
            })}
            onClick={() => handleExpandClick(index)}
          />
        </div>
        <div className={styles.moveInfoList}>{expandedMoves[index] && <MultiMovesMoveInfoList move={m} />}</div>
      </div>
    </React.Fragment>
  ));

  return (
    <div data-testid="move-container" className={styles.movesContainer}>
      {moveList}
    </div>
  );
};

export default MultiMovesMoveContainer;
