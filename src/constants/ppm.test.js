import { renderMultiplier } from './ppms';

describe('renderMultiplier', () => {
  it('null', () => {
    expect(renderMultiplier('')).toBe(null);
  });
  it('1.3', () => {
    expect(renderMultiplier(1.3)).toBe('(with 1.3x multiplier)');
  });
  it('1', () => {
    expect(renderMultiplier(1)).toBe('(with 1x multiplier)');
  });
});
