outputPath = System.getenv('DTC_OUTPUT_PATH') ?: 'build/docs'
inputPath  = 'project/req42-tracer/docs'

inputFiles = [
    [file: 'user-guide/index.adoc', formats: ['html', 'pdf']],
    [file: 'arc42/arc42.adoc',      formats: ['html']],
]

imageDirs = ["${inputPath}/images"]
