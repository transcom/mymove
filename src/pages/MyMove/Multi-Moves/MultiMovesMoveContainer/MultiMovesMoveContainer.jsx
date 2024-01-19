import React, { useState } from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import classnames from 'classnames';
import { Button } from '@trussworks/react-uswds';

import MultiMovesMoveInfoList from '../MultiMovesMoveInfoList/MultiMovesMoveInfoList';
import ButtonDropdownMenu from '../ButtonDropdownMenu/ButtonDropdownMenu';

import styles from './MultiMovesMoveContainer.module.scss';

const MultiMovesMoveContainer = ({ move }) => {
  const [expandedMoves, setExpandedMoves] = useState({});

  const handleExpandClick = (index) => {
    setExpandedMoves((prev) => ({
      ...prev,
      [index]: !prev[index],
    }));
  };

  const handleDropdownItemClick = (selectedItem) => {
    console.log(`${selectedItem.value}`);
  };

  const items = [
    {
      id: 1,
      value: 'PCS Orders',
    },
    {
      id: 2,
      value: 'PPM Packet',
    },
  ];

  const moveList = move.map((m, index) => (
    <React.Fragment key={index}>
      <div className={styles.moveContainer}>
        <div className={styles.heading} key={index}>
          <h3>#{m.moveCode}</h3>
          {m.status !== 'COMPLETED' ? (
            <Button className={styles.goToMoveBtn}>Go to Move</Button>
          ) : (
            <ButtonDropdownMenu
              title="Download"
              items={items}
              divClassName={styles.dropdownBtn}
              onItemClick={handleDropdownItemClick}
            />
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
