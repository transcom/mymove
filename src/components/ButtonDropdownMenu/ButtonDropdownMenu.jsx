import React, { useRef, useState } from 'react';
import { Button } from '@trussworks/react-uswds';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import classnames from 'classnames';

import styles from './ButtonDropdownMenu.module.scss';

function ButtonDropdownMenu({ title, items, multiSelect = false, divClassName, onItemClick, outline, minimal }) {
  const [open, setOpen] = useState(false);
  const [selection, setSelection] = useState([]);
  const toggle = () => setOpen(!open);
  const dropdownRef = useRef(null);

  const handleOutsideClick = (event) => {
    if (dropdownRef.current && !dropdownRef.current.contains(event.target)) {
      // Clicked outside the dropdown container, close the dropdown
      setOpen(false);
    }
  };

  const handleButtonClick = () => {
    toggle(!open);
    document.addEventListener('mousedown', handleOutsideClick);
  };

  function handleOnClick(item) {
    if (!selection.some((current) => current.id === item.id)) {
      if (!multiSelect) {
        setSelection([item]);
      } else if (multiSelect) {
        setSelection([...selection, item]);
      }
    } else {
      let selectionAfterRemoval = selection;
      selectionAfterRemoval = selectionAfterRemoval.filter((current) => current.id !== item.id);
      setSelection([...selectionAfterRemoval]);
    }

    // Call the onItemClick callback with the selected item
    if (onItemClick) {
      onItemClick(item);
      toggle(!open);
    }
  }

  let buttonContent;
  if (outline) {
    buttonContent = (
      <Button
        className={styles.btn}
        onClick={() => handleButtonClick()}
        secondary
        outline
        data-testid="button-dropdown-menu"
      >
        <span>{title}</span>
        <div>
          <FontAwesomeIcon icon="download" />
        </div>
      </Button>
    );
  } else if (minimal) {
    buttonContent = (
      <Button
        className={styles.btnMinimal}
        onClick={() => handleButtonClick()}
        secondary
        outline
        data-testid="button-dropdown-menu"
      >
        <span>{title}</span>
        <div>
          <FontAwesomeIcon icon="download" />
        </div>
      </Button>
    );
  } else {
    buttonContent = (
      <Button className={styles.btn} onClick={() => handleButtonClick()} data-testid="button-dropdown-menu">
        <span>{title}</span>
        <div>
          <FontAwesomeIcon icon="download" />
        </div>
      </Button>
    );
  }

  return (
    <div className={classnames(styles.dropdownWrapper, divClassName)}>
      <div className={styles.dropdownContainer} ref={dropdownRef}>
        {buttonContent}
        {open && (
          <ul className={styles.dropdownList}>
            {items.map((item) => (
              <li key={item.id} data-testid={`${item.id}-dropdown-item`}>
                <button type="button" onClick={() => handleOnClick(item)}>
                  <span>{item.value}</span>
                </button>
              </li>
            ))}
          </ul>
        )}
      </div>
    </div>
  );
}

export default ButtonDropdownMenu;
