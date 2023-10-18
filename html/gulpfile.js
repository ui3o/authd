const { src, dest, task } = require('gulp');
const clean = require('gulp-clean');
const inlineSource = require('gulp-inline-source');
const rename = require('gulp-rename');

task('clean', () => {
    return src('dist', { read: false, allowEmpty: true }).pipe(clean());
});

task('inline', () => {
    return src('dist/index.html').pipe(inlineSource()).pipe(rename('inline.html')).pipe(dest('dist/'));
});
