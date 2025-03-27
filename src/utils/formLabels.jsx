import React from 'react';

export const redAsterisk = <span style={{ color: 'red' }}>*</span>;

export const requiredAsteriskMessage = (
  <div data-testid="reqAsteriskMsg">Fields marked with {redAsterisk} are required.</div>
);

export const getLabelWithAsterisk = (label) => {
  return (
    <span data-testid="labelWithAsterisk">
      {label} {redAsterisk}
    </span>
  );
};
