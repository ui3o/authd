const { src, task } = require('gulp');
const clean = require('gulp-clean');
const inlineSource = require('gulp-inline-source');

task('clean', () => {
    return src('dist', { read: false, allowEmpty: true }).pipe(clean());
});

task('inline', () => {
    return src('dist/index.html').pipe(inlineSource());
});
